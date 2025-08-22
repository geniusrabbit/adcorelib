package actiontracker

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

// Option type
type Option[EventT EventType] func(ext *Extension[EventT])

// WithURLGenerator interface
func WithURLGenerator[EventT EventType](urlGenerator adtype.URLGenerator) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.urlGenerator = urlGenerator
	}
}

// WithEventStream setter
func WithEventStream[EventT EventType](eventStream eventstream.Stream) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.eventStream = eventStream
	}
}

// WithHTTPHandlerWrapper setter
func WithHTTPHandlerWrapper[EventT EventType](handlerWrapper *httphandler.HTTPHandlerWrapper) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.handlerWrapper = handlerWrapper
	}
}

// WithEventAllocator setter
func WithEventAllocator[EventT EventType](eventAllocator eventgenerator.Allocator[EventT]) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.eventAllocator = eventAllocator
	}
}

// WithCustomPriceExtractor setter
func WithCustomPriceExtractor[EventT EventType](priceExtractor PriceExtractor) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.priceExtractor = priceExtractor
	}
}

// WithDefaultPriceExtractor setter
func WithDefaultPriceExtractor[EventT EventType](paramName string) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.priceExtractor = DefaultPriceExtractor(paramName)
	}
}
