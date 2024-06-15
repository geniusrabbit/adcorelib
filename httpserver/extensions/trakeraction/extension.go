package trakeraction

import (
	"context"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/billing"
	"geniusrabbit.dev/adcorelib/context/ctxlogger"
	"geniusrabbit.dev/adcorelib/debugtool"
	"geniusrabbit.dev/adcorelib/eventtraking/events"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/fasttime"
	"geniusrabbit.dev/adcorelib/gtracing"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
	"geniusrabbit.dev/adcorelib/urlgenerator"
)

const (
	trakingGIFPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C" +
		"\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"
	trakingJSCode = "var __traking_time=new Date();"
)

// Extension of the server
type Extension struct {
	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// URL generator used for initialisation of submodules
	urlGenerator adtype.URLGenerator

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
	// Click/Direct links
	router.GET(ext.urlGenerator.ClickRouterURL(),
		ext.handlerWrapper.Metrics("click", ext.eventHandler(events.Click)))
	router.GET(ext.urlGenerator.DirectRouterURL(),
		ext.handlerWrapper.Metrics("direct", ext.eventHandler(events.Direct)))
	router.GET(ext.urlGenerator.WinRouterURL(),
		ext.handlerWrapper.Metrics("win", ext.eventHandler(events.AccessPointWin)))

	for _, leadType := range []string{"gif", "js", "plain"} {
		suffix := ""
		if leadType != "plain" {
			suffix = "." + leadType
		}
		router.GET("/lead"+suffix,
			ext.handlerWrapper.Metrics("traking.lead."+leadType, ext.eventLeadHandler(leadType)))
	}
}

// eventHandler by AD. This method works only for BaseSource
func (ext *Extension) eventHandler(eventName events.Type) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event." + eventName.String()
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		var (
			data    = rctx.QueryArgs().Peek("c")
			event   events.Event
			err     = event.Unpack(data, decodeEvents)
			span, _ = gtracing.StartSpanFromFastContext(rctx, handlerCode)
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil || eventName != event.Event {
			// TODO: Log event as unappropriate
			// s.Render.DirectBidResponse(nil, ctx)
			rctx.SetStatusCode(http.StatusBadRequest)
			ctxlogger.Get(ctx).Error("event handler", zap.Error(err),
				zap.String("event", event.Event.String()),
				zap.String("event_expected", eventName.String()),
			)
			return
		}

		if priceVal, _ := rctx.QueryArgs().GetUfloat("price"); priceVal > 0 {
			event.PurchaseViewPrice = billing.MoneyFloat(priceVal).Int64()
			if event.PurchaseViewPrice > event.ViewPrice {
				// Something went wrong
				ctxlogger.Get(ctx).Error("invalid win price, greater then original",
					zap.Float64("view_price", billing.Money(event.ViewPrice).Float64()),
					zap.Float64("purchase_view_price", billing.Money(event.PurchaseViewPrice).Float64()))
			}
		}

		if event.Event == events.Click || event.Event == events.Direct {
			rctx.Redirect(urlgenerator.PrepareURL(event.URL, &event), http.StatusFound)
		} else {
			rctx.SetStatusCode(http.StatusOK)
		}
		event.SetDateTime(int64(fasttime.UnixTimestampNano()))
		ext.eventStream.SendEvent(ctx, &event)
	}
}

func (ext *Extension) eventLeadHandler(pixelType string) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event.lead"
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		defer debugtool.Trace()

		var (
			dataCode  = []byte{}
			leadData  = rctx.QueryArgs().Peek("l")
			leadEvent events.LeadCode
			err       = leadEvent.Unpack(leadData)
			span, _   = gtracing.StartSpanFromFastContext(rctx, handlerCode)
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil {
			ctxlogger.Get(ctx).Error("event handler lead", zap.Error(err),
				zap.String("auction_id", leadEvent.AuctionID))
			rctx.SetStatusCode(http.StatusBadRequest)
			return
		} else {
			leadEvent.Timestamp = int64(fasttime.UnixTimestamp())
			ext.eventStream.SendLeadEvent(ctx, &leadEvent)
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
		rctx.Write(dataCode)
	}
}

func decodeEvents(code events.Code) events.Code {
	return code.URLDecode().Decompress()
}
