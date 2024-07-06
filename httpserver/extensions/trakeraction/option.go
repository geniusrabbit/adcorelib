package trakeraction

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
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
