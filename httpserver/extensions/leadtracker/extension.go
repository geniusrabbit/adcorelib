package leadtracker

import (
	"context"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/debugtool"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/gtracing"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

const (
	trakingGIFPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C" +
		"\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"
	trakingJSCode = "var __traking_time=new Date();"
)

type (
	LeadType = eventgenerator.LeadType
)

// Extension of the server
type Extension[LeadT LeadType] struct {
	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// URL generator used for initialisation of submodules
	urlGenerator adtype.URLGenerator

	// Event stream interface sends data into the queue
	eventStream eventstream.Stream

	// Lead allocator
	leadAllocator eventgenerator.Allocator[LeadT]
}

// NewExtension with options
func NewExtension[LeadT LeadType](opts ...Option[LeadT]) *Extension[LeadT] {
	ext := &Extension[LeadT]{}
	for _, opt := range opts {
		opt(ext)
	}
	return ext
}

// InitRouter of the HTTP server
func (ext *Extension[LeadType]) InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer) {
	for _, leadType := range []string{"gif", "js", "plain"} {
		suffix := ""
		if leadType != "plain" {
			suffix = "." + leadType
		}
		router.GET("/lead"+suffix,
			ext.handlerWrapper.Metrics("traking.lead."+leadType, ext.eventLeadHandler(leadType)))
	}
}

func (ext *Extension[LeadT]) eventLeadHandler(pixelType string) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event.lead"
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		defer debugtool.Trace()

		var (
			dataCode  []byte
			leadData  = rctx.QueryArgs().Peek("l")
			leadEvent = ext.leadAllocator()
			err       = leadEvent.Unpack(leadData)
			span, _   = gtracing.StartSpanFromFastContext(rctx, handlerCode)
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil {
			ctxlogger.Get(ctx).Error(
				"unpack event handler",
				zap.String("handler", "lead"),
				zap.String("auction_id", leadEvent.EventAuctionID()),
				zap.Error(err),
			)
			rctx.SetStatusCode(http.StatusBadRequest)
			return
		}

		// Send lead event to the stream
		leadEvent.SetDateTime(int64(fasttime.UnixTimestamp()))
		if err = ext.eventStream.SendLeadEvent(ctx, &leadEvent); err != nil {
			ctxlogger.Get(ctx).Error(
				"send event handler",
				zap.String("handler", "lead"),
				zap.String("auction_id", leadEvent.EventAuctionID()),
				zap.Error(err),
			)
		}

		if pixelType == "js" {
			rctx.SetContentType("application/javascript")
			dataCode = []byte(trakingJSCode)
		} else if pixelType == "gif" {
			rctx.SetContentType("image/gif")
			dataCode = []byte(trakingGIFPixel)
		} else {
			rctx.SetContentType("plain/text")
		}
		rctx.SetStatusCode(http.StatusOK)
		rctx.Response.Header.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		rctx.Response.Header.Set("Expires", "Wed, 11 Nov 1998 11:11:11 GMT")
		rctx.Response.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		rctx.Response.Header.Set("Pragma", "no-cache")
		_, _ = rctx.Write(dataCode)
	}
}
