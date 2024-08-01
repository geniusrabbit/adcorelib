package eventhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
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

// type metricsAccessor interface {
// 	Metrics() *openlatency.MetricsInfo
// }

// Extension of the server
type Extension struct {
	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// Event stream interface sends data into the queue
	eventStream eventstream.Stream
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
	// Pixel traking section
	router.GET("/t/px.gif",
		ext.handlerWrapper.Metrics("traking.pixel.gif", ext.eventSimpleHandler("gif")))
	router.GET("/t/px.js",
		ext.handlerWrapper.Metrics("traking.pixel.js", ext.eventSimpleHandler("js")))
}

func (ext *Extension) eventSimpleHandler(name string) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event." + name
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		var (
			pixelData = rctx.QueryArgs().Peek("i")
			event     events.Event
			err       = event.Unpack(pixelData, decodeEvents)
			span, _   = gtracing.StartSpanFromFastContext(rctx, handlerCode)
			dataCode  []byte
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil {
			ctxlogger.Get(ctx).Error("event handler "+name, zap.Error(err),
				zap.String("event", event.Event.String()))
			rctx.SetStatusCode(http.StatusBadRequest)
			return
		} else {
			event.SetDateTime(int64(fasttime.UnixTimestampNano()))
			// TODO: process error
			_ = ext.eventStream.SendEvent(ctx, &event)
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
