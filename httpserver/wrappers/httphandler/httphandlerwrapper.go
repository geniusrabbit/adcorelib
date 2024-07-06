package httphandler

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/net/fasthttp/middleware"
	"github.com/geniusrabbit/adcorelib/personification"
)

type (
	// ExtHTTPHandler with general context param
	ExtHTTPHandler func(context.Context, *fasthttp.RequestCtx)

	// ExtHTTPSpyHandler with general context param and person
	ExtHTTPSpyHandler func(context.Context, *fasthttp.RequestCtx, personification.Person)
)

// HTTPHandlerWrapper wraps extended HTTP handler to standart
type HTTPHandlerWrapper struct {
	// prepareRequestCtx params
	prepareRequestCtx func(ctx *fasthttp.RequestCtx)

	// newRequestCtx param creates new context for each request
	newRequestCtx func(ctx *fasthttp.RequestCtx) context.Context

	// Logger base object
	logger *zap.Logger
}

// NewHTTPHandlerWrapper returns new wrapper object
func NewHTTPHandlerWrapper(
	prepareRequestCtx func(ctx *fasthttp.RequestCtx),
	newRequestCtx func(ctx *fasthttp.RequestCtx) context.Context,
	logger *zap.Logger,
) *HTTPHandlerWrapper {
	return &HTTPHandlerWrapper{
		prepareRequestCtx: prepareRequestCtx,
		newRequestCtx:     newRequestCtx,
		logger:            logger,
	}
}

// Metrics wraps default Ext HTTP handler with metrics
func (wrp *HTTPHandlerWrapper) Metrics(name string, handler ExtHTTPHandler) fasthttp.RequestHandler {
	newHandler := wrp.HTTPHandler(handler)
	handlerWrapper := middleware.CollectSimpleMetrics(name, newHandler)
	if wrp.prepareRequestCtx == nil {
		return middleware.Logger(wrp.logger, handlerWrapper)
	}
	return middleware.Logger(wrp.logger, func(ctx *fasthttp.RequestCtx) {
		wrp.prepareRequestCtx(ctx)
		handlerWrapper(ctx)
	})
}

// HTTPHandler wraps default Ext HTTP handler
func (wrp *HTTPHandlerWrapper) HTTPHandler(f ExtHTTPHandler) fasthttp.RequestHandler {
	return func(req *fasthttp.RequestCtx) {
		ctx := wrp.newRequestContext(req)
		f(ctx, req)
	}
}

// SpyMetrics wraps default Ext HTTP handler with metrics and personification
func (wrp *HTTPHandlerWrapper) SpyMetrics(name string, spy middleware.Spy, handler ExtHTTPSpyHandler) fasthttp.RequestHandler {
	newHandler := wrp.HTTPSpyHandler(spy, handler)
	handlerWrapper := middleware.CollectSimpleMetrics(name, newHandler)
	if wrp.prepareRequestCtx == nil {
		return middleware.Logger(wrp.logger, handlerWrapper)
	}
	return middleware.Logger(wrp.logger, func(ctx *fasthttp.RequestCtx) {
		wrp.prepareRequestCtx(ctx)
		handlerWrapper(ctx)
	})
}

// HTTPSpyHandler wraps default Ext HTTP handler and personification
func (wrp *HTTPHandlerWrapper) HTTPSpyHandler(spy middleware.Spy, f ExtHTTPSpyHandler) fasthttp.RequestHandler {
	return spy(func(person personification.Person, req *fasthttp.RequestCtx) {
		ctx := wrp.newRequestContext(req)
		f(ctx, req, person)
	})
}

func (wrp *HTTPHandlerWrapper) newRequestContext(req *fasthttp.RequestCtx) context.Context {
	if wrp.newRequestCtx != nil {
		return wrp.newRequestCtx(req)
	}
	return context.Background()
}
