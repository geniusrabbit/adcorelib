package leadtracker

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

// Option type
type Option[LeadT LeadType] func(ext *Extension[LeadT])

// WithURLGenerator interface
func WithURLGenerator[LeadT LeadType](urlGenerator adtype.URLGenerator) Option[LeadT] {
	return func(ext *Extension[LeadT]) {
		ext.urlGenerator = urlGenerator
	}
}

// WithEventStream setter
func WithEventStream[LeadT LeadType](eventStream eventstream.Stream) Option[LeadT] {
	return func(ext *Extension[LeadT]) {
		ext.eventStream = eventStream
	}
}

// WithHTTPHandlerWrapper setter
func WithHTTPHandlerWrapper[LeadT LeadType](handlerWrapper *httphandler.HTTPHandlerWrapper) Option[LeadT] {
	return func(ext *Extension[LeadT]) {
		ext.handlerWrapper = handlerWrapper
	}
}

// WithLeadAllocator setter
func WithLeadAllocator[LeadT LeadType](leadAllocator eventgenerator.Allocator[LeadT]) Option[LeadT] {
	return func(ext *Extension[LeadT]) {
		ext.leadAllocator = leadAllocator
	}
}
