package adsourceexperiments

import (
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"
)

type sourceSimpleWrapper struct {
	source adtype.Source
}

// NewSimpleWrapper object
func NewSimpleWrapper(source adtype.Source) SourceWrapper {
	return &sourceSimpleWrapper{source: source}
}

// Next returns source interface according to strategy
func (w *sourceSimpleWrapper) Next() adtype.Source {
	return w.source
}

// SetTimeout for sourcer
func (w *sourceSimpleWrapper) SetTimeout(timeout time.Duration) {
	if src, _ := w.source.(adtype.SourceTimeoutSetter); src != nil {
		src.SetTimeout(timeout)
	}
}

var _ SourceWrapper = (*sourceSimpleWrapper)(nil)
