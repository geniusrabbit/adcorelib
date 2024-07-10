package adsource

import (
	"fmt"
	"time"

	"github.com/geniusrabbit/adcorelib/adsource/experiments"
	"github.com/geniusrabbit/adcorelib/adtype"
)

// Option sets some property of the server
type Option func(wrp *MultisourceWrapper)

// WithBaseSource as default
func WithBaseSource(source any) Option {
	return func(wrp *MultisourceWrapper) {
		switch src := source.(type) {
		case nil:
		case adtype.Source:
			wrp.baseSource = experiments.NewSimpleWrapper(src)
		case experiments.SourceWrapper:
			wrp.baseSource = src
		default:
			panic(fmt.Sprintf("Invalid base source type %T", source))
		}
	}
}

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
