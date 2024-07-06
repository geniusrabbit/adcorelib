//
// @project geniusrabbit::sspserver 2017, 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2024
//

package adsource

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/demdxx/rpool/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/adsourceexperiments"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/auction"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/gtracing"
)

// Error set...
var (
	ErrSourcesCantBeNil = errors.New("[SSP] seurces cant be nil")
)

const (
	minimalTimeout          = time.Millisecond * 10
	minimalParallelRequests = 1
)

// MultisourceWrapper describes the abstraction which can control what request
// where should be sent in which driver
type MultisourceWrapper struct {
	// Main source which called everytime
	baseSource adsourceexperiments.SourceWrapper

	// Source list of external platforms
	sources adtype.SourceAccessor

	// Execution pool
	execpool *rpool.Pool

	// RequestTimeout duration
	requestTimeout time.Duration

	// MaxParallelRequest number
	maxParallelRequest int

	// Metrics accessor
	metrics Metrics
}

// NewMultisourceWrapper SSP inited with options
func NewMultisourceWrapper(options ...Option) (*MultisourceWrapper, error) {
	var wrp MultisourceWrapper

	for _, opt := range options {
		opt(&wrp)
	}

	if wrp.sources == nil {
		return nil, ErrSourcesCantBeNil
	}

	if wrp.requestTimeout < minimalTimeout {
		wrp.requestTimeout = minimalTimeout
	}

	if wrp.maxParallelRequest < minimalParallelRequests {
		wrp.maxParallelRequest = minimalParallelRequests
	}

	wrp.execpool = rpool.NewPool(rpool.WithMaxTasksCount(wrp.maxParallelRequest))

	return &wrp, nil
}

// ID of the source driver
func (wrp *MultisourceWrapper) ID() uint64 { return 0 }

// Protocol of the source driver
func (wrp *MultisourceWrapper) Protocol() string { return "multisource" }

// Test request before processing
func (wrp *MultisourceWrapper) Test(request *adtype.BidRequest) bool { return true }

// Bid request for standart system filter
func (wrp *MultisourceWrapper) Bid(request *adtype.BidRequest) (response adtype.Responser) {
	if wrp == nil {
		return adtype.NewEmptyResponse(request, nil, errors.New("wrapper is nil"))
	}
	var (
		err     error
		count   = wrp.maxParallelRequest
		tube    = make(chan adtype.Responser, wrp.maxParallelRequest)
		span, _ = gtracing.StartSpanFromContext(request.Ctx, "ssp.bid")
		referee auction.Referee
		timeout bool
	)

	if span != nil {
		ext.Component.Set(span, "ssp")
		oldContext := request.Ctx
		request.Ctx = opentracing.ContextWithSpan(oldContext, span)
		defer func() {
			request.Ctx = oldContext
			span.Finish()
		}()
	}

	// Base request to internal DB
	if src := wrp.getMainSource(); wrp != nil && wrp.testSource(src, request) {
		startTime := fasttime.UnixTimestampNano()
		response := src.Bid(request)
		wrp.metrics.IncrementBidRequestCount(src, request, time.Duration(startTime-fasttime.UnixTimestampNano()))

		// Store bidding information
		wrp.sourceResponseLog(request, response)

		if response.Error() == nil {
			referee.Push(response.Ads()...)
			// TODO update minimal bids by response
			// TODO release response
		} else {
			wrp.metrics.IncrementBidErrorCount(src, request, response.Error())
		}
	}

	// Source request loop
	for iterator := wrp.sources.Iterator(request); ; {
		src := iterator.Next()
		if src == nil {
			break
		}

		if wrp.testSource(src, request) {
			count--
			wrp.execpool.Go(func() {
				startTime := fasttime.UnixTimestampNano()
				response := src.Bid(request)
				wrp.metrics.IncrementBidRequestCount(src, request, time.Duration(startTime-fasttime.UnixTimestampNano()))
				tube <- response

				// Store bidding information
				wrp.sourceResponseLog(request, response)

				if response.Error() != nil {
					wrp.metrics.IncrementBidErrorCount(src, request, response.Error())
				}
			})
			if src.RequestStrategy().IsSingle() {
				break
			}
		}

		if count < 1 {
			break
		}
	}

	// Auction loop
	if count < wrp.maxParallelRequest {
		timer := time.NewTimer(wrp.requestTimeout)
		for ; count < wrp.maxParallelRequest; count++ {
			select {
			case resp := <-tube:
				if respErr := resp.Error(); respErr != nil {
					err = respErr
				} else {
					referee.Push(resp.Ads()...)
				}
			case <-timer.C:
				count = wrp.maxParallelRequest
				timeout = true
			}
		}

		if !timeout {
			timer.Stop()
		}
	}

	if items := referee.MatchRequest(request); len(items) > 0 {
		response = adtype.BorrowResponse(request, nil, items, nil)
		err = nil
	} else {
		response = adtype.NewEmptyResponse(request, wrp, err)
	}

	return response
}

// ProcessResponse when need to fix the result and process all counters
func (wrp *MultisourceWrapper) ProcessResponse(response adtype.Responser) {
	if response == nil || response.Error() != nil {
		return
	}
	// Pricess prices of campaigns
	for _, it := range response.Ads() {
		if it.Validate() != nil {
			continue
		}
		switch ad := it.(type) {
		case adtype.ResponserItem:
			wrp.processAdResponse(response, ad)
		case adtype.ResponserMultipleItem:
			for _, it := range ad.Ads() {
				wrp.processAdResponse(response, it)
			}
		default:
			ctxlogger.Get(response.Context()).
				Warn("Unsupportable respont item type", zap.String("type", fmt.Sprintf("%T", it)))
		}
	}
}

// ProcessResponseItem result or error
func (wrp *MultisourceWrapper) ProcessResponseItem(response adtype.Responser, item adtype.ResponserItem) {
	if src := item.Source(); src != nil {
		src.ProcessResponseItem(response, item)
	}
}

// SetRequestTimeout of the simple request
func (wrp *MultisourceWrapper) SetRequestTimeout(timeout time.Duration) {
	if timeout < minimalTimeout {
		timeout = minimalTimeout
	}
	if wrp.requestTimeout != timeout {
		wrp.requestTimeout = timeout
		wrp.sources.SetTimeout(timeout)
		if wrp.baseSource != nil {
			wrp.baseSource.SetTimeout(timeout)
		}
	}
}

// Sources of the ads
func (wrp *MultisourceWrapper) Sources() adtype.SourceAccessor {
	return wrp.sources
}

// RequestStrategy description
func (wrp *MultisourceWrapper) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}

// RevenueShareReduceFactor which is a potential
func (wrp *MultisourceWrapper) RevenueShareReduceFactor() float64 { return 0 }

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (wrp *MultisourceWrapper) PriceCorrectionReduceFactor() float64 { return 0 }

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (wrp *MultisourceWrapper) sourceResponseLog( /* bidRequest */ _ *adtype.BidRequest, response adtype.Responser) {
	if isNil(response) {
		return
	}

	eventStream := eventstream.StreamFromContext(response.Context())
	if response.Error() == nil && len(response.Ads()) > 0 {
		// Log ads
		for _, it := range response.Ads() {
			switch ad := it.(type) {
			case adtype.ResponserItem:
				_ = eventStream.Send(events.SourceBid, events.StatusSuccess, response, ad)
			case adtype.ResponserMultipleItem:
				if len(ad.Ads()) > 0 {
					_ = eventStream.Send(events.SourceBid, events.StatusSuccess, response, ad.Ads()[0])
				}
			}
			break
		}
	} else if response.Error() == nil && len(response.Ads()) == 0 {
		_ = eventStream.SendSourceNoBid(response)
	} else if response.Error() != nil && strings.Contains(response.Error().Error(), "skip") {
		_ = eventStream.SendSourceSkip(response)
	} else {
		_ = eventStream.SendSourceFail(response)
	}
}

func (wrp *MultisourceWrapper) processAdResponse(response adtype.Responser, ad adtype.ResponserItem) {
	if src := ad.Source(); src != nil {
		src.ProcessResponseItem(response, ad)
	}
}

func (wrp *MultisourceWrapper) testSource(src adtype.Source, request *adtype.BidRequest) bool {
	return src != nil && request.SourceFilterCheck(src.ID()) && src.Test(request)
}

func (wrp *MultisourceWrapper) getMainSource() adtype.Source {
	if wrp.baseSource == nil {
		return nil
	}
	return wrp.baseSource.Next()
}

func isNil(v any) bool {
	switch v.(type) {
	case nil:
		return true
	}
	return false
}
