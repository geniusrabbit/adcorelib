package adtype

import (
	"time"
)

// SourceMinimal contains only minimal set of methods
type SourceMinimal interface {
	// Bid request for standart system filter
	Bid(request *BidRequest) Responser

	// ProcessResponseItem result or error
	ProcessResponseItem(Responser, ResponserItem)
}

// Source of advertisement and where will be selled the traffic
type Source interface {
	SourceMinimal

	// ID of the source driver
	ID() uint64

	// ObjectKey of the source driver
	ObjectKey() uint64

	// Protocol of the source driver
	Protocol() string

	// Test request before processing
	Test(request *BidRequest) bool

	// PriceCorrectionReduceFactor which is a potential
	// Returns percent from 0 to 1 for reducing of the value
	// If there is 10% of price correction, it means that 10% of the final price must be ignored
	PriceCorrectionReduceFactor() float64

	// RequestStrategy description
	RequestStrategy() RequestStrategy
}

// SourceTesteChecker checker
type SourceTesteChecker interface {
	// Test current request for compatibility
	Test(request *BidRequest) bool
}

// SourceTimeoutSetter interface
type SourceTimeoutSetter interface {
	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}

// SourceTester interface
type SourceTester interface {
	Source
	SourceTesteChecker
}
