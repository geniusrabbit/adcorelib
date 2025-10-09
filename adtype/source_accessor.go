package adtype

import (
	"context"
	"iter"
	"time"
)

// SourceIterator is a source iterator
type SourceIterator = iter.Seq2[float32, Source]

// SourceAccessor preoritise the source access
type SourceAccessor interface {
	// Iterator returns the configured queue accessor
	Iterator(request BidRequester) SourceIterator

	// SourceByID returns source instance
	SourceByID(ctx context.Context, id uint64) (Source, error)

	// SetTimeout for sourcer
	SetTimeout(ctx context.Context, timeout time.Duration)
}
