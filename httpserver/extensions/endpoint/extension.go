package endpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/xtypes"
	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httptraceroute"
	"github.com/geniusrabbit/adcorelib/net/fasthttp/middleware"
	"github.com/geniusrabbit/adcorelib/openlatency"
	"github.com/geniusrabbit/adcorelib/personification"
)

type (
	appAccessor interface {
		AppByURI(uri string) (*admodels.Application, error)
	}
	zoneAccessor interface {
		TargetByCodename(string) (adtype.Target, error)
	}
	getSourceAccessor interface {
		Sources() adtype.SourceAccessor
	}
	factoryListAccessor interface {
		FactoryList() []adtype.SourceFactory
	}
	sourceListAccessor interface {
		SourceList() ([]adtype.Source, error)
	}
	sourceMetricsAccessor interface {
		Metrics() *openlatency.MetricsInfo
	}
)

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

	// App data accessor
	appAccessor appAccessor

	// Zone data accessor
	zoneAccessor zoneAccessor

	// URL query pattern like `/b/{endpoint}/{zone}` by default
	URLQueryPattern string

	// List of endpoints of classic executors
	endpoints []Endpoint

	// Metrics
	adRequestCountMetrics *prometheus.CounterVec
}

// NewExtension with options
func NewExtension(opts ...Option) *Extension {
	ext := &Extension{
		adRequestCountMetrics: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "ad_request_count",
			Help: "Count of Ad requests",
		}, []string{"endpoint", "zone", "adblock", "robot", "secure", "private", "proxy"}),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(ext)
		}
	}

	return ext
}

// InitRouter of the HTTP server
func (ext *Extension) InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer) {
	routeWrapper := httptraceroute.Wrap(router, tracer)
	urlPattern := gocast.Or(ext.URLQueryPattern, "/b/{endpoint}/{zone}")

	// Ad handler request
	for _, endpoint := range ext.endpoints {
		endpointCode := endpoint.Codename()
		pattern := strings.ReplaceAll(urlPattern, "{endpoint}", endpointCode)

		routeWrapper.GET(pattern,
			ext.handlerWrapper.SpyMetrics("endpoint."+endpointCode, ext.spy,
				// Double wrap to evoid potential `endpoint` relink
				func(endpoint Endpoint) httphandler.ExtHTTPSpyHandler {
					return func(ctx context.Context, req *fasthttp.RequestCtx, person personification.Person) {
						ext.endpointRequestHandler(ctx, req, person, endpoint)
					}
				}(endpoint),
			))
	}

	// API info handlers
	if sa, ok := ext.source.(getSourceAccessor); ok {
		sources := sa.Sources()
		// Factories info API handler
		if fa, ok := sources.(factoryListAccessor); ok {
			routeWrapper.GET("/protocols", middleware.CollectSimpleMetrics("api.protocols", ext.factoryListHandler(fa)))
			for _, factory := range fa.FactoryList() {
				if factory.Info().Protocol == "" {
					ctxlogger.Get(ctx).Warn("Empty protocol in factory info", zap.String("factory", factory.Info().Name))
					continue
				}
				routeWrapper.GET("/protocols/"+factory.Info().Protocol,
					middleware.CollectSimpleMetrics("api.protocols."+factory.Info().Protocol, ext.factoryInfoHandler(factory)))
			}
		}

		// Source info API handler
		if sla, ok := sources.(sourceListAccessor); ok {
			routeWrapper.GET("/sources", middleware.CollectSimpleMetrics("api.sources", ext.sourceListHandler(sla)))
		}
		routeWrapper.GET("/sources/{id}",
			middleware.CollectSimpleMetrics("api.sources.info", ext.sourceInfoHandler(sources)))
		routeWrapper.GET("/sources/{id}/metrics",
			middleware.CollectSimpleMetrics("api.sources.metrics", ext.sourceMetricsHandler(sources)))
	}
}

func (ext *Extension) endpointRequestHandler(ctx context.Context, req *fasthttp.RequestCtx, person personification.Person, endpoint Endpoint) {
	bidRequest := ext.requestByHTTPRequest(ctx, person, req)

	if bidRequest == nil {
		// Collect metrics
		ext.adRequestCountMetrics.WithLabelValues(
			endpoint.Codename(),
			strings.Split(req.UserValue("zone").(string), ".")[0],
			"", "", "", "", "",
		).Inc()

		req.SetStatusCode(http.StatusNotFound)
		return
	}

	// Collect metrics
	ext.adRequestCountMetrics.WithLabelValues(
		endpoint.Codename(),
		bidRequest.Imps[0].Target.Codename(),
		b2sbool(bidRequest.IsAdblock()),
		b2sbool(bidRequest.IsRobot()),
		b2sbool(bidRequest.IsSecure()),
		b2sbool(bidRequest.IsPrivateBrowsing()),
		b2sbool(bidRequest.IsProxy()),
	).Inc()

	response := endpoint.Handle(ext.source, bidRequest)
	ext.source.ProcessResponse(response)
}

func (ext *Extension) factoryListHandler(fa factoryListAccessor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		factories := fa.FactoryList()
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_ = json.NewEncoder(ctx).Encode(map[string]any{
			"protocols": xtypes.SliceApply(factories, func(f adtype.SourceFactory) any {
				info := f.Info()
				return map[string]any{
					"protocol":    info.Protocol,
					"name":        info.Name,
					"description": info.Description,
					"versions":    info.Versions,
				}
			}),
		})
	}
}

func (ext *Extension) factoryInfoHandler(factory adtype.SourceFactory) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		info := factory.Info()
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_ = json.NewEncoder(ctx).Encode(info)
	}
}

func (ext *Extension) sourceListHandler(sa sourceListAccessor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		sources, err := sa.SourceList()
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_ = json.NewEncoder(ctx).Encode(xtypes.SliceApply(sources, func(s adtype.Source) any {
			return map[string]any{
				"id":       s.ID(),
				"protocol": s.Protocol(),
			}
		}))
	}
}

func (ext *Extension) sourceInfoHandler(sa adtype.SourceAccessor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		id := gocast.Uint64(ctx.UserValue("id"))
		src, _ := sa.SourceByID(id)
		if src == nil {
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			_, _ = ctx.Write([]byte(`{"error":"source not found"}`))
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_ = json.NewEncoder(ctx).Encode(map[string]any{
			"id":       src.ID(),
			"protocol": src.Protocol(),
		})
	}
}

func (ext *Extension) sourceMetricsHandler(sa adtype.SourceAccessor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		id := gocast.Uint64(ctx.UserValue("id"))
		src, _ := sa.SourceByID(id)
		if src == nil {
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetContentType("application/json")
			_, _ = ctx.Write([]byte(`{"error":"source not found"}`))
			return
		}
		var metrics any
		if sm, ok := src.(sourceMetricsAccessor); ok {
			metrics = sm.Metrics()
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		_ = json.NewEncoder(ctx).Encode(metrics)
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func (ext *Extension) requestByHTTPRequest(ctx context.Context, person personification.Person, rctx *fasthttp.RequestCtx) *adtype.BidRequest {
	var (
		app       *admodels.Application
		spotInfo  = strings.Split(rctx.UserValue("zone").(string), ".")
		target, _ = ext.zoneAccessor.TargetByCodename(spotInfo[0])
	)
	if target == nil {
		return nil
	}

	// Get application by referer
	if ext.appAccessor != nil {
		app, _ = ext.appAccessor.AppByURI(domain(string(rctx.Referer())))
	}

	// Create request for the target
	return NewRequestFor(ctx, app, target, person,
		NewRequestOptions(rctx), ext.formatAccessor)
}
