package pixeltracker

import (
	"context"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
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

// EventType object for event basic type interface
type EventType = eventgenerator.EventType

// Extension of the server
type Extension[EventT EventType] struct {
	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// Event stream interface sends data into the queue
	eventStream eventstream.Stream

	// Event allocator
	eventAllocator eventgenerator.Allocator[EventT]
}

// NewExtension with options
func NewExtension[EventT EventType](opts ...Option[EventT]) *Extension[EventT] {
	ext := &Extension[EventT]{}
	for _, opt := range opts {
		opt(ext)
	}
	return ext
}

// InitRouter of the HTTP server
func (ext *Extension[EventT]) InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer) {
	// Pixel traking section
	router.GET("/t/px.gif",
		ext.handlerWrapper.Metrics("traking.pixel.gif", ext.eventSimpleHandler("gif")))
	router.GET("/t/px.js",
		ext.handlerWrapper.Metrics("traking.pixel.js", ext.eventSimpleHandler("js")))
}

func (ext *Extension[EventT]) eventSimpleHandler(name string) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event." + name
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		var (
			pixelData = rctx.QueryArgs().Peek("i")
			event     = ext.eventAllocator()
			err       = event.Unpack(pixelData, decodeEvents)
			span, _   = gtracing.StartSpanFromFastContext(rctx, handlerCode)
			dataCode  []byte
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil {
			ctxlogger.Get(ctx).Error("unpack event handler",
				zap.String("handler", name),
				zap.String("event", event.EventType().String()),
				zap.Error(err),
			)
		} else {
			event.SetDateTime(int64(fasttime.UnixTimestampNano()))
			if err = ext.eventStream.SendEvent(ctx, &event); err != nil {
				ctxlogger.Get(ctx).Error("send event handler",
					zap.String("handler", name),
					zap.String("event", event.EventType().String()),
					zap.Error(err),
				)
			}
		}

		switch name {
		case "js":
			rctx.SetContentType("application/javascript")
			dataCode = []byte(trakingJSCode)
		case "gif":
			rctx.SetContentType("image/gif")
			dataCode = []byte(trakingGIFPixel)
		}
		rctx.SetStatusCode(http.StatusOK)
		rctx.Response.Header.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		rctx.Response.Header.Set("Expires", "Wed, 11 Nov 1998 11:11:11 GMT")
		rctx.Response.Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		rctx.Response.Header.Set("Pragma", "no-cache")
		_, _ = rctx.Write(dataCode)
	}
}

func decodeEvents(code events.Code) events.Code {
	return code.URLDecode().Decompress()
}
