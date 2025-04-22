package adtype

import (
	"iter"
	"time"
)

// SourceIterator is a source iterator
type SourceIterator = iter.Seq2[float32, Source]

// SourceAccessor preoritise the source access
type SourceAccessor interface {
	// Iterator returns the configured queue accessor
	Iterator(request *BidRequest) SourceIterator

	// SourceByID returns source instance
	SourceByID(id uint64) (Source, error)

	// SetTimeout for sourcer
	SetTimeout(timeout time.Duration)
}
