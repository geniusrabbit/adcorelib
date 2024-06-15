package pixel

import (
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
)

// Option type
type Option func(ext *Extension)

// WithEventStream setter
func WithEventStream(eventStream eventstream.Stream) Option {
	return func(ext *Extension) {
		ext.eventStream = eventStream
	}
}

// WithHTTPHandlerWrapper setter
func WithHTTPHandlerWrapper(handlerWrapper *httphandler.HTTPHandlerWrapper) Option {
	return func(ext *Extension) {
		ext.handlerWrapper = handlerWrapper
	}
}
