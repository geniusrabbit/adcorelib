package adtype

import (
	"time"
)

// RequestStrategy defines politics of request sending
type RequestStrategy int

const (
	// AsynchronousRequestStrategy is default strategy implies
	// requesting all auction participants and choising the most
	// profitable variant of all
	AsynchronousRequestStrategy RequestStrategy = iota

	// SingleRequestStrategy tells that if response was
	// received it should be performed
	SingleRequestStrategy
)

func (rs RequestStrategy) IsSingle() bool {
	return rs == SingleRequestStrategy
}

func (rs RequestStrategy) IsAsynchronous() bool {
	return rs == AsynchronousRequestStrategy
}

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

// SourceIterator returns next source from the scope
type SourceIterator interface {
	// Next returns source interface according to strategy
	Next() Source
}

// SourceAccessor preoritise the source access
type SourceAccessor interface {
	// Iterator returns the configured queue accessor
	Iterator(request *BidRequest) SourceIterator

	// SourceByID returns source instance
	SourceByID(id uint64) (Source, error)

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}
