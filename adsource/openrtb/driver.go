package openrtb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bsm/openrtb"
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/adsource/optimizer"
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/context/ctxlogger"
	counter "geniusrabbit.dev/adcorelib/errorcounter"
	"geniusrabbit.dev/adcorelib/eventtraking/events"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/fasttime"
	"geniusrabbit.dev/adcorelib/openlatency"
)

const (
	headerRequestOpenRTBVersion  = "X-Openrtb-Version"
	headerRequestOpenRTBVersion2 = "2.5"
	headerRequestOpenRTBVersion3 = "3.0"
	defaultMinWeight             = 0.001
)

type driver struct {
	lastRequestTime uint64

	// Requests RPS counter
	rpsCurrent     counter.Counter
	errorCounter   counter.ErrorCounter
	latencyMetrics *openlatency.MetricsCounter

	// Original source model
	source *admodels.RTBSource

	// Request headers
	Headers map[string]string

	// Client of HTTP requests
	Client *http.Client

	// optimizator object
	optimizator *optimizer.Optimizer
}

func newDriver(_ context.Context, source *admodels.RTBSource, _ ...any) (*driver, error) {
	var headers map[string]string
	if source.Headers.Data != nil {
		headers = *source.Headers.Data
	}
	if source.MinimalWeight <= 0 {
		source.MinimalWeight = defaultMinWeight
	}
	return &driver{
		source:         source,
		latencyMetrics: &openlatency.MetricsCounter{},
		Headers:        headers,
		optimizator:    optimizer.New(),
	}, nil
}

// ID of source
func (d *driver) ID() uint64 { return d.source.ID }

// Test request before processing
func (d *driver) Test(request *adtype.BidRequest) bool {
	if d.source.RPS > 0 {
		if !d.source.Options.ErrorsIgnore && !d.errorCounter.Next() {
			return false
		}

		now := fasttime.UnixTimestampNano()
		if now-atomic.LoadUint64(&d.lastRequestTime) >= uint64(time.Second) {
			atomic.StoreUint64(&d.lastRequestTime, now)
			d.rpsCurrent.Set(0)
		} else if d.rpsCurrent.Get() >= int64(d.source.RPS) {
			return false
		}
	}

	if !d.source.Test(request) {
		return false
	}

	// Check formats targeting
	for _, f := range request.Formats() {
		if !d.optimizator.Test(
			uint(f.ID),
			byte(request.GeoID()),
			request.LanguageID(),
			request.DeviceID(),
			request.OSID(),
			request.BrowserID(),
			d.source.MinimalWeight,
		) {
			return false
		}
	}
	return true
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (d *driver) PriceCorrectionReduceFactor() float64 {
	return d.source.PriceCorrectionReduceFactor()
}

// RequestStrategy description
func (d *driver) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}

// Bid request for standart system filter
func (d *driver) Bid(request *adtype.BidRequest) (response adtype.Responser) {
	d.rpsCurrent.Inc(1)

	httpRequest, err := d.request(request)
	if err != nil {
		return adtype.NewErrorResponse(request, err)
	}

	resp, err := d.getClient().Do(httpRequest)
	if err != nil {
		d.processHTTPReponse(resp, err)
		ctxlogger.Get(request.Ctx).Debug("bid",
			zap.String("source_url", d.source.URL),
			zap.Error(err))
		return adtype.NewErrorResponse(request, err)
	}

	ctxlogger.Get(request.Ctx).Debug("bid",
		zap.String("source_url", d.source.URL),
		zap.String("http_response_status_txt", http.StatusText(resp.StatusCode)),
		zap.Int("http_response_status", resp.StatusCode))

	if resp.StatusCode == http.StatusNoContent {
		return adtype.NewErrorResponse(request, ErrNoCampaignsStatus)
	}

	if resp.StatusCode != http.StatusOK {
		d.processHTTPReponse(resp, nil)
		return adtype.NewErrorResponse(request, ErrInvalidResponseStatus)
	}

	defer resp.Body.Close()
	if res, err := d.unmarshal(request, resp.Body); d.source.Options.Trace && err != nil {
		response = adtype.NewErrorResponse(request, err)
		ctxlogger.Get(request.Ctx).Error("bid response", zap.Error(err))
	} else if res != nil {
		response = res
	}

	if response != nil && response.Error() == nil {
		for _, ad := range response.Ads() {
			for _, f := range ad.Impression().Formats() {
				_ = d.optimizator.Inc(
					uint(f.ID),
					byte(request.GeoID()),
					request.LanguageID(),
					request.DeviceID(),
					request.OSID(),
					request.BrowserID(),
					d.source.MinimalWeight,
				)
			}
		}
	} else {
		for _, f := range request.Formats() {
			_ = d.optimizator.Inc(
				uint(f.ID),
				byte(request.GeoID()),
				request.LanguageID(),
				request.DeviceID(),
				request.OSID(),
				request.BrowserID(),
				-d.source.MinimalWeight,
			)
		}
	}

	d.processHTTPReponse(resp, err)
	if response == nil {
		return adtype.NewEmptyResponse(request, d, err)
	}
	return response
}

// ProcessResponseItem result or error
func (d *driver) ProcessResponseItem(response adtype.Responser, item adtype.ResponserItem) {
	if response == nil || response.Error() != nil {
		return
	}
	for _, ad := range response.Ads() {
		switch bid := ad.(type) {
		case *adtype.ResponseBidItem:
			if len(bid.Bid.NURL) > 0 {
				eventstream.WinsFromContext(response.Context()).Send(response.Context(), bid.Bid.NURL)
				ctxlogger.Get(response.Context()).Info("ping", zap.String("url", bid.Bid.NURL))
			}
			eventstream.StreamFromContext(response.Context()).
				Send(events.Impression, events.StatusUndefined, response, bid)
		default:
			// Dummy...
		}
	}
}

// RevenueShareReduceFactor which is a potential
func (d *driver) RevenueShareReduceFactor() float64 {
	return 0 // TODO: d.source.PriceCorrectionReduce / 100
}

///////////////////////////////////////////////////////////////////////////////
/// Implementation of platform.Metrics interface
///////////////////////////////////////////////////////////////////////////////

// Metrics information of the platform
func (d *driver) Metrics() *openlatency.MetricsInfo {
	var info openlatency.MetricsInfo
	d.latencyMetrics.FillMetrics(&info)
	info.ID = d.ID()
	info.Protocol = d.source.Protocol
	return &info
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (d *driver) getClient() *http.Client {
	if d.Client == nil {
		timeout := time.Millisecond * time.Duration(d.source.Timeout)
		if timeout < 1 {
			timeout = time.Millisecond * 150
		}
		d.Client = &http.Client{Timeout: timeout}
	}
	return d.Client
}

// prepare request for RTB
func (d *driver) request(request *adtype.BidRequest) (req *http.Request, err error) {
	var (
		rtbRequest interface{}
		data       bytes.Buffer
	)

	if d.source.Protocol == "openrtb3" {
		rtbRequest = requestToRTBv3(request, d.getRequestOptions()...)
	} else {
		rtbRequest = requestToRTBv2(request, d.getRequestOptions()...)
	}

	// Prepare data for request
	if err = json.NewEncoder(&data).Encode(rtbRequest); err != nil {
		return nil, err
	}

	// Create new request
	if req, err = http.NewRequest(d.source.Method, d.source.URL, &data); err != nil {
		return nil, err
	}

	d.fillRequest(request, req)
	return req, nil
}

func (d *driver) unmarshal(request *adtype.BidRequest, r io.Reader) (response *adtype.BidResponse, err error) {
	var bidResp openrtb.BidResponse

	switch d.source.RequestType {
	case RequestTypeJSON:
		if d.source.Options.Trace {
			var data []byte
			if data, err = io.ReadAll(r); err == nil {
				ctxlogger.Get(request.Ctx).Error("unmarshal",
					zap.String("src_url", d.source.URL),
					zap.String("unmarshal_data", string(data)))
				err = json.Unmarshal(data, &bidResp)
			}
		} else {
			err = json.NewDecoder(r).Decode(&bidResp)
		}
	case RequestTypeXML, RequestTypeProtobuff:
		err = fmt.Errorf("request body type not supported: %s", d.source.RequestType.Name())
	default:
		err = fmt.Errorf("undefined request type: %s", d.source.RequestType.Name())
	}

	if err != nil {
		return nil, err
	}

	// Check response for support HTTPS
	if request.IsSecure() {
		for _, seat := range bidResp.SeatBid {
			for _, bid := range seat.Bid {
				if strings.Contains(bid.AdMarkup, "http://") {
					return nil, ErrResponseAreNotSecure
				}
			}
		} // end for
	}

	// Build response
	response = &adtype.BidResponse{
		Src:         d,
		Req:         request,
		BidResponse: bidResp,
	}
	response.Prepare()
	return response, nil
}

// fillRequest of HTTP
func (d *driver) fillRequest(request *adtype.BidRequest, httpReq *http.Request) {
	httpReq.Header.Set("Content-Type", "application/json")
	if d.source.Protocol == "openrtb3" {
		httpReq.Header.Set(headerRequestOpenRTBVersion, headerRequestOpenRTBVersion3)
	} else {
		httpReq.Header.Set(headerRequestOpenRTBVersion, headerRequestOpenRTBVersion2)
	}
	httpReq.Header.Set(openlatency.HTTPHeaderRequestTimemark, strconv.FormatInt(openlatency.RequestInitTime(request.Time()), 10))

	// Fill default headers
	for key, value := range d.Headers {
		httpReq.Header.Set(key, value)
	}
}

// @link https://golang.org/src/net/http/status.go
func (d *driver) processHTTPReponse(resp *http.Response, err error) {
	switch {
	case err != nil || resp == nil ||
		(resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent):
		// if err == http.ErrHandlerTimeout {

		// }
		d.errorCounter.Inc()
	default:
		d.errorCounter.Dec()
	}
}

func (d *driver) getRequestOptions() []BidRequestRTBOption {
	return []BidRequestRTBOption{
		WithRTBOpenNativeVersion("1.1"),
		WithFormatFilter(d.source.TestFormat),
		WithMaxTimeDuration(time.Duration(d.source.Timeout) * time.Millisecond),
	}
}
