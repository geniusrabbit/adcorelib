package adsourceexperiments

import (
	"time"

	"geniusrabbit.dev/adcorelib/adtype"
)

// SourceWrapper advertisement accessor interface
type SourceWrapper interface {
	// Next returns source interface according to strategy
	Next() adtype.Source

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}
