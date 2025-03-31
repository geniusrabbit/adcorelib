package adtype

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/platform/info"
)

// SourceFactory is a source factory interface
type SourceFactory interface {
	// Info returns information about the source platform
	// and the source protocol
	Info() info.Platform

	// Protocols returns list of supported protocols of the source
	Protocols() []string

	// New creates a new source instance for the given RTBSource
	New(ctx context.Context, source *admodels.RTBSource, opts ...any) (SourceTester, error)
}
