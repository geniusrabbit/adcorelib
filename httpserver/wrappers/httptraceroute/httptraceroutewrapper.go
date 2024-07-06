package httptraceroute

import (
	"github.com/fasthttp/router"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/gtracing"
)

// TraceRouterWrapper of the HTTP router
type TraceRouterWrapper struct {
	router *router.Router
	tracer opentracing.Tracer
}

// Wrap router
func Wrap(router *router.Router, tracer opentracing.Tracer) *TraceRouterWrapper {
	return &TraceRouterWrapper{
		router: router,
		tracer: tracer,
	}
}

// GET HTTP action
func (w TraceRouterWrapper) GET(path string, h fasthttp.RequestHandler) {
	w.router.GET(path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}

// POST HTTP action
func (w TraceRouterWrapper) POST(path string, h fasthttp.RequestHandler) {
	w.router.POST(path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}

// Handler HTTP action
func (w TraceRouterWrapper) Handle(method, path string, h fasthttp.RequestHandler) {
	w.router.Handle(method, path, gtracing.FastHTTPTraceWrapper(w.tracer, h))
}
