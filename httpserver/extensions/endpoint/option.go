package endpoint

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
	"github.com/geniusrabbit/adcorelib/net/fasthttp/middleware"
)

// Option type
type Option func(ext *Extension)

// WithAdvertisementSource accessor
func WithAdvertisementSource(source Source) Option {
	return func(ext *Extension) {
		ext.source = source
	}
}

// WithHTTPHandlerWrapper setter
func WithHTTPHandlerWrapper(handlerWrapper *httphandler.HTTPHandlerWrapper) Option {
	return func(ext *Extension) {
		ext.handlerWrapper = handlerWrapper
	}
}

// WithFormatAccessor setter
func WithFormatAccessor(formatAccessor types.FormatsAccessor) Option {
	return func(ext *Extension) {
		ext.formatAccessor = formatAccessor
	}
}

// WithAppAccessor setter
func WithAppAccessor(appAccessor appAccessor) Option {
	return func(ext *Extension) {
		ext.appAccessor = appAccessor
	}
}

// WithZoneAccessor setter
func WithZoneAccessor(zoneAccessor zoneAccessor) Option {
	return func(ext *Extension) {
		ext.zoneAccessor = zoneAccessor
	}
}

// WithSpy setter
func WithSpy(spy middleware.Spy) Option {
	return func(ext *Extension) {
		ext.spy = spy
	}
}

// WithSendpoints setter
func WithSendpoints(endpoints ...Endpoint) Option {
	return func(ext *Extension) {
		ext.endpoints = ext.endpoints[:0]
		for _, endpoint := range endpoints {
			if endpoint != nil {
				ext.endpoints = append(ext.endpoints, endpoint)
			}
		}
	}
}
