// Package adformattest provides test fixtures for adformat — the
// counterpart of admodels/types.MockFormats(), but backed by embed.FS
// (via adformat/examples) instead of the old, fragile
// runtime.Caller(1) + os.Open("../assets/...") pair, which broke when a
// package was moved or vendored.
package adformattest

import (
	"github.com/geniusrabbit/adcorelib/adformat"
	"github.com/geniusrabbit/adcorelib/adformat/examples"
)

// MockFormats returns one *adformat.Format per embedded
// adformat/examples/*.yaml sample (banner, banner_fixed, popunder,
// native, video_preroll, push_ad, playable, multiscreen), decoded and
// Validate()-checked once at call time. IDs are reassigned sequentially
// (1..N) so the returned slice can be fed straight into
// adformat.NewSimpleAccessor without collisions — two of the embedded
// examples intentionally reuse id: 1 in their own YAML file, since each
// is written as a standalone, independent example.
//
// Panics on decode/validation failure, same convention as the old
// MockFormats() — these are compile-time-known fixtures, a failure here
// means adformat/examples itself is broken, not bad runtime input.
func MockFormats() []*adformat.Format {
	formats, err := examples.All()
	if err != nil {
		panic(err)
	}
	for i, f := range formats {
		f.ID = uint64(i + 1)
	}
	return formats
}

// NewAccessor builds a ready-to-use adformat.SimpleAccessor over
// MockFormats(), for tests that need a working Accessor rather than a
// bare format list.
func NewAccessor() *adformat.SimpleAccessor {
	return adformat.NewSimpleAccessor(MockFormats()...)
}

// FormatByCodename returns the mock format with the given codename, or
// nil. A small convenience on top of MockFormats() for tests that only
// care about one specific example.
func FormatByCodename(codename string) *adformat.Format {
	for _, f := range MockFormats() {
		if f.Codename == codename {
			return f
		}
	}
	return nil
}
