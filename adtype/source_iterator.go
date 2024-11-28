package adtype

import "time"

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
