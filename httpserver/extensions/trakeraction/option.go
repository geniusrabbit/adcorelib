package trakeraction

import (
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
)

// Option type
type Option func(ext *Extension)

// WithURLGenerator interface
func WithURLGenerator(urlGenerator adtype.URLGenerator) Option {
	return func(ext *Extension) {
		ext.urlGenerator = urlGenerator
	}
}

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
