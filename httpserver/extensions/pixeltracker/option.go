package pixeltracker

import (
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

// Option type
type Option[EventT EventType] func(ext *Extension[EventT])

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

func WithEventAllocator[EventT EventType](eventAllocator eventgenerator.Allocator[EventT]) Option[EventT] {
	return func(ext *Extension[EventT]) {
		ext.eventAllocator = eventAllocator
	}
}
