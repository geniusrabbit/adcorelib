package openrtb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bsm/openrtb"
	"github.com/demdxx/gocast/v2"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	counter "github.com/geniusrabbit/adcorelib/errorcounter"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/net/httpclient"
	"github.com/geniusrabbit/adcorelib/openlatency"
	"github.com/geniusrabbit/adcorelib/openlatency/prometheuswrapper"
)

const (
	headerRequestOpenRTBVersion  = "X-Openrtb-Version"
	headerRequestOpenRTBVersion2 = "2.5"
	headerRequestOpenRTBVersion3 = "3.0"
	defaultMinWeight             = 0.001
)

type driver[NetDriver httpclient.Driver[Rq, Rs], Rq httpclient.Request, Rs httpclient.Response] struct {
	lastRequestTime uint64

	// Requests RPS counter
	rpsCurrent     counter.Counter
	errorCounter   counter.ErrorCounter
	latencyMetrics *prometheuswrapper.Wrapper

	// Original source model
	source *admodels.RTBSource

	// Request headers
	Headers map[string]string

	// Client of HTTP requests
	netClient NetDriver
}

func newDriver[ND httpclient.Driver[Rq, Rs], Rq httpclient.Request, Rs httpclient.Response](_ context.Context, source *admodels.RTBSource, netClient ND, _ ...any) (*driver[ND, Rq, Rs], error) {
	if source.MinimalWeight <= 0 {
		source.MinimalWeight = defaultMinWeight
	}
	return &driver[ND, Rq, Rs]{
		source:    source,
		Headers:   source.Headers.DataOr(nil),
		netClient: netClient,
		latencyMetrics: prometheuswrapper.NewWrapperDefault("adsource_",
			[]string{"id", "protocol", "driver"},
			[]string{gocast.Str(source.ID), source.Protocol, "openrtb"},
		),
	}, nil
}

// ID of source
func (d *driver[ND, Rq, Rs]) ID() uint64 { return d.source.ID }

// Protocol of source
func (d *driver[ND, Rq, Rs]) Protocol() string { return d.source.Protocol }

// Test request before processing
func (d *driver[ND, Rq, Rs]) Test(request *adtype.BidRequest) bool {
	if d.source.RPS > 0 {
		if !d.source.Options.ErrorsIgnore && !d.errorCounter.Next() {
			d.latencyMetrics.IncSkip()
			return false
		}

		now := fasttime.UnixTimestampNano()
		if now-atomic.LoadUint64(&d.lastRequestTime) >= uint64(time.Second) {
			atomic.StoreUint64(&d.lastRequestTime, now)
			d.rpsCurrent.Set(0)
		} else if d.rpsCurrent.Get() >= int64(d.source.RPS) {
			d.latencyMetrics.IncSkip()
			return false
		}
	}

	if !d.source.Test(request) {
		d.latencyMetrics.IncSkip()
		return false
	}

	return true
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (d *driver[ND, Rq, Rs]) PriceCorrectionReduceFactor() float64 {
	return d.source.PriceCorrectionReduceFactor()
}

// RequestStrategy description
func (d *driver[ND, Rq, Rs]) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}

// Bid request for standart system filter
func (d *driver[ND, Rq, Rs]) Bid(request *adtype.BidRequest) (response adtype.Responser) {
	beginTime := time.Now().UnixNano()
	d.rpsCurrent.Inc(1)
	d.latencyMetrics.BeginQuery()

	httpRequest, err := d.request(request)
	if err != nil {
		return adtype.NewErrorResponse(request, err)
	}

	resp, err := d.netClient.Do(httpRequest)
	d.latencyMetrics.UpdateQueryLatency(time.Duration(time.Now().UnixNano() - beginTime))

	if err != nil {
		d.processHTTPReponse(resp, err)
		ctxlogger.Get(request.Ctx).Debug("bid",
			zap.String("source_url", d.source.URL),
			zap.Error(err))
		return adtype.NewErrorResponse(request, err)
	}

	ctxlogger.Get(request.Ctx).Debug("bid",
		zap.String("source_url", d.source.URL),
		zap.String("http_response_status_txt", http.StatusText(resp.StatusCode())),
		zap.Int("http_response_status", resp.StatusCode()))

	if resp.StatusCode() == http.StatusNoContent {
		d.latencyMetrics.IncNobid()
		return adtype.NewErrorResponse(request, ErrNoCampaignsStatus)
	}

	if resp.StatusCode() != http.StatusOK {
		d.processHTTPReponse(resp, nil)
		return adtype.NewErrorResponse(request, ErrInvalidResponseStatus)
	}

	defer resp.Close()
	if res, err := d.unmarshal(request, resp.Body()); d.source.Options.Trace && err != nil {
		response = adtype.NewErrorResponse(request, err)
		ctxlogger.Get(request.Ctx).Error("bid response", zap.Error(err))
	} else if res != nil {
		response = res
	}

	if response != nil && response.Error() == nil {
		if len(response.Ads()) > 0 {
			d.latencyMetrics.IncSuccess()
		} else {
			d.latencyMetrics.IncNobid()
		}
	}

	d.processHTTPReponse(resp, err)
	if response == nil {
		return adtype.NewEmptyResponse(request, d, err)
	}
	return response
}

// ProcessResponseItem result or error
func (d *driver[ND, Rq, Rs]) ProcessResponseItem(response adtype.Responser, item adtype.ResponserItem) {
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
func (d *driver[ND, Rq, Rs]) RevenueShareReduceFactor() float64 {
	return 0 // TODO: d.source.PriceCorrectionReduce / 100
}

///////////////////////////////////////////////////////////////////////////////
/// Implementation of platform.Metrics interface
///////////////////////////////////////////////////////////////////////////////

// Metrics information of the platform
func (d *driver[ND, Rq, Rs]) Metrics() *openlatency.MetricsInfo {
	var info openlatency.MetricsInfo
	d.latencyMetrics.FillMetrics(&info)
	info.ID = d.ID()
	info.Protocol = d.source.Protocol
	info.QPSLimit = d.source.RPS
	return &info
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

// prepare request for RTB
func (d *driver[ND, Rq, Rs]) request(request *adtype.BidRequest) (req Rq, err error) {
	var (
		rtbRequest any
		data       bytes.Buffer
	)

	if d.source.Protocol == "openrtb3" {
		rtbRequest = requestToRTBv3(request, d.getRequestOptions()...)
	} else {
		rtbRequest = requestToRTBv2(request, d.getRequestOptions()...)
	}

	// Prepare data for request
	if err = json.NewEncoder(&data).Encode(rtbRequest); err != nil {
		return d.netClient.NoopRequest(), err
	}

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	// enc.Encode(rtbRequest)

	// Create new request
	if req, err = d.netClient.Request(d.source.Method, d.source.URL, &data); err != nil {
		return req, err
	}

	d.fillRequest(request, req)
	return req, nil
}

func (d *driver[ND, Rq, Rs]) unmarshal(request *adtype.BidRequest, r io.Reader) (response *adtype.BidResponse, err error) {
	var bidResp openrtb.BidResponse

	switch d.source.RequestType {
	case RequestTypeJSON:
		if d.source.Options.Trace {
			var data []byte
			if data, err = io.ReadAll(r); err == nil {
				ctxlogger.Get(request.Ctx).Error("trace unmarshal",
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
func (d *driver[ND, Rq, Rs]) fillRequest(request *adtype.BidRequest, httpReq Rq) {
	httpReq.SetHeader("Content-Type", "application/json")
	if d.source.Protocol == "openrtb3" {
		httpReq.SetHeader(headerRequestOpenRTBVersion, headerRequestOpenRTBVersion3)
	} else {
		httpReq.SetHeader(headerRequestOpenRTBVersion, headerRequestOpenRTBVersion2)
	}
	httpReq.SetHeader(openlatency.HTTPHeaderRequestTimemark,
		strconv.FormatInt(openlatency.RequestInitTime(request.Time()), 10))

	// Fill default headers
	for key, value := range d.Headers {
		httpReq.SetHeader(key, value)
	}
}

// @link https://golang.org/src/net/http/status.go
func (d *driver[ND, Rq, Rs]) processHTTPReponse(resp Rs, err error) {
	switch {
	case err != nil || resp.IsNoop() ||
		(resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent):
		if errors.Is(err, http.ErrHandlerTimeout) {
			d.latencyMetrics.IncTimeout()
		}
		d.errorCounter.Inc()
		d.latencyMetrics.IncError(openlatency.MetricErrorHTTP, http.StatusText(resp.StatusCode()))
	default:
		d.errorCounter.Dec()
	}
}

func (d *driver[ND, Rq, Rs]) getRequestOptions() []BidRequestRTBOption {
	return []BidRequestRTBOption{
		WithRTBOpenNativeVersion("1.1"),
		WithFormatFilter(d.source.TestFormat),
		WithMaxTimeDuration(time.Duration(d.source.Timeout) * time.Millisecond),
	}
}
