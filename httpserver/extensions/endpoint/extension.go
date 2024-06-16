package endpoint

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/sspserver/udetect"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httptraceroute"
	"geniusrabbit.dev/adcorelib/net/fasthttp/middleware"
	"geniusrabbit.dev/adcorelib/personification"
)

type zoneAccessor interface {
	TargetByID(uint64) (admodels.Target, error)
}

// Extension of the server
type Extension struct {
	// Source of the advertisement
	source Source

	// Spy wrapper
	spy middleware.Spy

	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// Format accessor
	formatAccessor types.FormatsAccessor

	// Zone data accessor
	zoneAccessor zoneAccessor

	// List of endpoints of classic executors
	endpoints []Endpoint
}

// NewExtension with options
func NewExtension(opts ...Option) *Extension {
	ext := &Extension{}
	for _, opt := range opts {
		opt(ext)
	}
	return ext
}

// InitRouter of the HTTP server
func (ext *Extension) InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer) {
	routeWrapper := httptraceroute.Wrap(router, tracer)

	// Ad handler request
	for _, endpoint := range ext.endpoints {
		routeWrapper.GET("/b/"+endpoint.Codename()+"/{zone}",
			ext.handlerWrapper.SpyMetrics("endpoint."+endpoint.Codename(), ext.spy,
				// Double wrap to evoid potential `endpoint` relink
				func(endpoint Endpoint) httphandler.ExtHTTPSpyHandler {
					return func(ctx context.Context, req *fasthttp.RequestCtx, person personification.Person) {
						ext.endpointRequestHandler(ctx, req, person, endpoint)
					}
				}(endpoint),
			))
	}
}

func (ext *Extension) endpointRequestHandler(ctx context.Context, req *fasthttp.RequestCtx, person personification.Person, endpoint Endpoint) {
	bidRequest := ext.requestByHTTPRequest(ctx, person, req)
	if bidRequest == nil {
		req.SetStatusCode(http.StatusNotFound)
		return
	}
	response := endpoint.Handle(ext.source, bidRequest)
	ext.source.ProcessResponse(response)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (ext *Extension) requestByHTTPRequest(ctx context.Context, person personification.Person, rctx *fasthttp.RequestCtx) *adtype.BidRequest {
	var (
		spotInfo  = strings.Split(rctx.UserValue("zone").(string), ".")
		zoneID, _ = strconv.ParseUint(spotInfo[0], 10, 64)
		target, _ = ext.zoneAccessor.TargetByID(zoneID)
	)
	if target == nil {
		return nil
	}
	request := NewRequestFor(ctx, target, person, NewRequestOptions(rctx), ext.formatAccessor)
	if request != nil {
		ext.prepareRequest(request)
	}
	return request
}

func (ext *Extension) prepareRequest(request *adtype.BidRequest) {
	// BigBrother, tell me! Who is it?
	var (
		query    = request.RequestCtx.QueryArgs()
		keywords = peekOneFromQuery(query, "keywords", "keyword", "kw")
	)

	if request.Person != nil {
		// Fill user info
		ui := request.Person.UserInfo()
		request.User = &adtype.User{
			ID:            ui.UUID(),
			SessionID:     ui.SessionID(),
			FingerPrintID: ui.Fingerprint(),
			ETag:          ui.ETag(),
			Keywords:      keywords,
			Geo:           ui.GeoInfo(),
		}
		request.Device = ui.DeviceInfo()
		if ui != nil && ui.User != nil {
			request.User.AgeStart = ui.User.AgeStart
			request.User.AgeEnd = ui.User.AgeEnd
		}
	} else {
		request.User = &adtype.User{
			FingerPrintID: "",
			Keywords:      keywords,
		}
	}

	// Fill GEO info
	if request.User.Geo == nil {
		request.User.Geo = &udetect.Geo{
			IP: request.RequestCtx.RemoteIP(),
		}
	}
}
