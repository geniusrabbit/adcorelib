package actiontracker

import (
	"context"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/gtracing"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

type (
	EventType = eventgenerator.EventType
)

// Extension of the server
type Extension[EventT EventType] struct {
	// Price extractor function
	priceExtractor PriceExtractor

	// Wrapper of extended handler to default
	handlerWrapper *httphandler.HTTPHandlerWrapper

	// URL generator used for initialisation of submodules
	urlGenerator adtype.URLGenerator

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
	// Click/Direct links
	router.GET(ext.urlGenerator.ClickRouterURL(),
		ext.handlerWrapper.Metrics("click", ext.eventHandler(events.Click)))
	router.GET(ext.urlGenerator.DirectRouterURL(),
		ext.handlerWrapper.Metrics("direct", ext.eventHandler(events.Direct)))
	router.GET(ext.urlGenerator.WinRouterURL(),
		ext.handlerWrapper.Metrics("win", ext.eventHandler(events.AccessPointWin)))
}

// eventHandler by AD. This method works only for BaseSource
func (ext *Extension[EventT]) eventHandler(eventName events.Type) httphandler.ExtHTTPHandler {
	handlerCode := "postback.event." + eventName.String()
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) {
		var (
			data    = rctx.QueryArgs().Peek("c")
			event   = ext.eventAllocator()
			err     = event.Unpack(data, decodeEvents)
			span, _ = gtracing.StartSpanFromFastContext(rctx, handlerCode)
		)

		if span != nil {
			defer span.Finish()
		}

		if err != nil || eventName != event.EventType() {
			// TODO: Log event as unappropriate
			// s.Render.DirectBidResponse(nil, ctx)
			rctx.SetStatusCode(http.StatusBadRequest)
			ctxlogger.Get(ctx).Error("event handler",
				zap.String("handler", eventName.String()),
				zap.String("event", event.EventType().String()),
				zap.Error(err),
			)
			return
		}

		// Set custom price for the event
		if ext.priceExtractor != nil {
			if priceVal, _ := ext.priceExtractor(ctx, rctx); priceVal > 0 {
				err := event.SetEventPurchaseViewPrice(billing.MoneyFloat(priceVal).Int64())
				if err != nil {
					// Something went wrong
					ctxlogger.Get(ctx).Error("new purchase price",
						zap.String("handler", eventName.String()),
						zap.String("event", event.EventType().String()),
						zap.Error(err),
					)
				}
			}
		}

		// Redirect user to the target URL if it's needed for the event type
		if etype := event.EventType(); etype == events.Click || etype == events.Direct {
			rctx.Redirect(event.PrepareURL(event.EventURL()), http.StatusFound)
		} else {
			rctx.SetStatusCode(http.StatusOK)
		}

		// Send action event to the stream
		event.SetDateTime(int64(fasttime.UnixTimestampNano()))
		if err = ext.eventStream.SendEvent(ctx, &event); err != nil {
			ctxlogger.Get(ctx).Error("send event handler",
				zap.String("handler", eventName.String()),
				zap.String("event", event.EventType().String()),
				zap.Error(err),
			)
		}
	}
}

func decodeEvents(code events.Code) events.Code {
	return code.URLDecode().Decompress()
}
