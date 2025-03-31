package adsource

import (
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// Option sets some property of the server
type Option func(wrp *MultisourceWrapper)

// WithSourceAccessor for the server
func WithSourceAccessor(sources adtype.SourceAccessor) Option {
	return func(wrp *MultisourceWrapper) {
		wrp.sources = sources
	}
}

// WithTimeout of one request
func WithTimeout(timeout time.Duration) Option {
	return func(wrp *MultisourceWrapper) {
		wrp.requestTimeout = timeout
	}
}

// WithMaxParallelRequests returns count of requests to external sources by one request
func WithMaxParallelRequests(maxParallelRequest int) Option {
	return func(wrp *MultisourceWrapper) {
		wrp.maxParallelRequest = maxParallelRequest
	}
}
