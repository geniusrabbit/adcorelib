// Package adsource provides the implementation of ad source drivers for the AdEngine.
// This package offers a unified interface for interacting with various ad sources
// and includes a collection of standard ad sources such as in-memory, database,
// OpenRTB (Real-Time Bidding) sources, and more.
//
// The primary components of the package are:
// - MultisourceWrapper: An abstraction that manages multiple ad sources and controls
//   the request distribution and response collection from these sources.
// - Error definitions: Standardized error messages used across the package.
// - Internal methods: Helper functions and internal logic to support the primary
//   operations of the MultisourceWrapper and other components.
//
// The package ensures efficient and parallel processing of ad requests by utilizing
// a worker pool for executing bid requests. It also integrates with tracing and
// logging systems to provide detailed insights into the bidding process and performance.
//
// Key Features:
// - Unified interface for multiple ad sources
// - Support for parallel bid requests
// - Integration with tracing (using OpenTracing) and logging (using zap)
// - Metrics collection and reporting for monitoring performance
//
// Example usage:
//   wrapper, err := adsource.NewMultisourceWrapper(options...)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   response := wrapper.Bid(request)
//   if response.Error() != nil {
//       log.Println("Bid request failed:", response.Error())
//   } else {
//       log.Println("Bid request succeeded:", response.Ads())
//   }

package adsource

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/demdxx/rpool/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/geniusrabbit/adcorelib/adquery/bidresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/auction/trafaret"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/gtracing"
)

// Error set...
var (
	ErrSourcesCantBeNil = errors.New("[SSP] seurces can`t be nil")
)

const (
	minimalTimeout          = time.Millisecond * 10
	minimalParallelRequests = 1
)

type respItem struct {
	priority float32
	resp     adtype.Response
}

// MultisourceWrapper describes the abstraction which can control where to send requests
// and how to handle responses from different sources.
type MultisourceWrapper struct {
	// Source list of external platforms
	sources adtype.SourceAccessor

	// Execution pool
	execpool *rpool.Pool

	// Request timeout duration
	requestTimeout time.Duration

	// Maximum number of parallel requests
	maxParallelRequest int

	// Metrics accessor
	metrics Metrics
}

// NewMultisourceWrapper initializes a new MultisourceWrapper with the given options
func NewMultisourceWrapper(options ...Option) (*MultisourceWrapper, error) {
	wrp := new(MultisourceWrapper)

	for _, opt := range options {
		opt(wrp)
	}

	if wrp.sources == nil {
		return nil, ErrSourcesCantBeNil
	}

	wrp.requestTimeout = max(wrp.requestTimeout, minimalTimeout)
	wrp.maxParallelRequest = max(wrp.maxParallelRequest, minimalParallelRequests)
	wrp.execpool = rpool.NewPool(rpool.WithMaxTasksCount(wrp.maxParallelRequest))

	return wrp, nil
}

// ID returns the ID of the source driver
func (wrp *MultisourceWrapper) ID() uint64 { return 0 }

// ObjectKey returns the object key of the source driver
func (wrp *MultisourceWrapper) ObjectKey() uint64 { return 0 }

// Protocol returns the protocol of the source driver
func (wrp *MultisourceWrapper) Protocol() string { return "multisource" }

// Test validates the request before processing
func (wrp *MultisourceWrapper) Test(request adtype.BidRequester) bool { return true }

// Bid handles a bid request and processes it through the appropriate sources
func (wrp *MultisourceWrapper) Bid(request adtype.BidRequester) (response adtype.Response) {
	if wrp == nil {
		return bidresponse.NewEmptyResponse(request, nil, errors.New("wrapper is nil"))
	}
	var (
		count         = wrp.maxParallelRequest
		isQueueClosed atomic.Bool
		queue         = make(chan respItem, count)
		span, _       = gtracing.StartSpanFromContext(request.Context(), "ssp.bid")
		trafaret      trafaret.Filler
		err           error
	)

	if span != nil {
		ext.Component.Set(span, "ssp")
		oldContext := request.Context()
		request.SetContext(opentracing.ContextWithSpan(oldContext, span))
		defer func() {
			request.SetContext(oldContext)
			span.Finish()
		}()
	}

	// Ensure that the queue is closed when the function exits
	defer func() {
		isQueueClosed.Store(true)
		close(queue)
	}()

	// Source request loop
	for prior, src := range wrp.sources.Iterator(request) {
		count--
		if isQueueClosed.Load() {
			break
		}
		wrp.execpool.Go(func() {
			if isQueueClosed.Load() {
				return
			}

			startTime := fasttime.UnixTimestampNano()

			// Send request to the source for the advertising
			resp := src.Bid(request)

			// Update metrics
			wrp.metrics.IncrementBidRequestCount(src,
				request, time.Duration(startTime-fasttime.UnixTimestampNano()))

			if !isQueueClosed.Load() {
				// Send response to the channel if it is still open
				select {
				case queue <- respItem{priority: prior, resp: resp}:
					// Successfully sent to the channel
				default:
					// Channel is closed or full, skip sending
				}
			}

			// Store bidding information
			wrp.sourceResponseLog(request, resp)

			if resp.Error() != nil {
				wrp.metrics.IncrementBidErrorCount(src, request, resp.Error())
			}
		})
		if src.RequestStrategy().IsSingle() || count < 1 {
			break
		}
	}

	// Auction loop processing with timeout
	if count < wrp.maxParallelRequest {
		timer := time.NewTimer(wrp.requestTimeout)
		defer func() {
			if !timer.Stop() {
				select {
				case <-timer.C: // Drain the channel if the timer already fired
				default:
				}
			}
		}()

		for ; count < wrp.maxParallelRequest; count++ {
			select {
			case item := <-queue:
				if respErr := item.resp.Error(); respErr != nil {
					err = respErr
				} else {
					trafaret.Push(item.priority, item.resp.Ads()...)
				}
			case <-timer.C:
				count = wrp.maxParallelRequest
			case <-request.Done():
				count = wrp.maxParallelRequest
			}
		}
	}

	// Prepare response
	{
		var items []adtype.ResponseItemCommon
		for _, imp := range request.Impressions() {
			if impItems := trafaret.Fill(imp.ID, imp.Count); len(impItems) > 0 {
				items = append(items, impItems...)
			}
		}

		if len(items) == 0 {
			response = bidresponse.NewEmptyResponse(request, wrp, err)
		} else {
			response = bidresponse.BorrowResponse(request, nil, items, nil)
		}
	}

	return response
}

// ProcessResponse processes the response to update metrics and log information
func (wrp *MultisourceWrapper) ProcessResponse(response adtype.Response) {
	if response == nil || response.Error() != nil {
		return
	}
	// Process prices of campaigns
	for ad := range response.IterAds() {
		if ad != nil && ad.Validate() == nil {
			wrp.ProcessResponseItem(response, ad)
		}
	}
}

// ProcessResponseItem processes an individual response item
func (wrp *MultisourceWrapper) ProcessResponseItem(response adtype.Response, ad adtype.ResponseItem) {
	if src := ad.Source(); src != nil {
		src.ProcessResponseItem(response, ad)
	}
}

// SetRequestTimeout sets the request timeout, ensuring it is not below the minimal timeout
func (wrp *MultisourceWrapper) SetRequestTimeout(ctx context.Context, timeout time.Duration) {
	if timeout < minimalTimeout {
		timeout = minimalTimeout
	}
	if wrp.requestTimeout != timeout {
		wrp.requestTimeout = timeout
		wrp.sources.SetTimeout(ctx, timeout)
	}
}

// Sources returns the source accessor
func (wrp *MultisourceWrapper) Sources() adtype.SourceAccessor {
	return wrp.sources
}

// RequestStrategy returns the request strategy
func (wrp *MultisourceWrapper) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}

// RevenueShareReduceFactor returns the revenue share reduce factor
func (wrp *MultisourceWrapper) RevenueShareReduceFactor() float64 { return 0 }

// PriceCorrectionReduceFactor returns the price correction reduce factor
// If there is a 10% price correction, it means that 10% of the final price must be ignored
func (wrp *MultisourceWrapper) PriceCorrectionReduceFactor() float64 { return 0 }

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (wrp *MultisourceWrapper) sourceResponseLog( /* bidRequest */ _ adtype.BidRequester, response adtype.Response) {
	if isNil(response) {
		return
	}

	eventStream := eventstream.StreamFromContext(response.Context())
	respErr := response.Error()

	// Check if response is valid and has ads items
	if respErr == nil && len(response.Ads()) > 0 {
		// Send bid events for each ad separately
		for ad := range response.IterAds() {
			eventType := events.Undefined
			eventStatus := events.StatusSuccess
			// Check ad response item
			if err := ad.Validate(); err != nil {
				if errors.Is(err, adtype.ErrResponseItemEmpty) {
					eventType = events.SourceNoBid
				} else if errors.Is(err, adtype.ErrResponseItemSkipped) {
					eventType = events.SourceSkip
				} else {
					eventType = events.SourceFail
					eventStatus = events.StatusFailed
				}
			} else {
				eventType = events.SourceBid
			}

			// Send event to event stream
			_ = eventStream.Send(eventType, uint8(eventStatus), response, ad)
		}

		// Send no bid for each empty slot (zone, adunit)
		imps := response.Request().Impressions()
		for _, imp := range imps {
			if response.Item(imp.ID) == nil {
				_ = eventStream.Send(events.SourceNoBid, events.StatusUndefined, response,
					&adtype.ResponseItemEmpty{Req: response.Request(), Imp: imp})
			}
		}
	} else if respErr == nil && len(response.Ads()) == 0 {
		_ = eventStream.SendSourceNoBid(response)
	} else if respErr != nil && (errors.Is(respErr, adtype.ErrResponseSkipped) || strings.Contains(respErr.Error(), "skip")) {
		_ = eventStream.SendSourceSkip(response)
	} else {
		_ = eventStream.SendSourceFail(response)
	}
}

//go:inline
func isNil(v any) bool {
	switch vv := v.(type) {
	case nil:
		return true
	case any:
		return vv == nil
	}
	return false
}
