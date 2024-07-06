package adsourceexperiments

import (
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// SourceWrapper advertisement accessor interface
type SourceWrapper interface {
	// Next returns source interface according to strategy
	Next() adtype.Source

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}
