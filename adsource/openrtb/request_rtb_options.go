package openrtb

import (
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// BidRequestRTBOptions of request build
type BidRequestRTBOptions struct {
	OpenNative struct {
		Ver string
	}
	FormatFilter func(f *types.Format) bool
	Currency     []string
	TimeMax      time.Duration
}

func (opts *BidRequestRTBOptions) openNativeVer() string {
	return opts.OpenNative.Ver
}

func (opts *BidRequestRTBOptions) currencies() []string {
	if len(opts.Currency) > 0 {
		return opts.Currency
	}
	return []string{"USD"}
}

// BidRequestRTBOption set function
type BidRequestRTBOption func(opts *BidRequestRTBOptions)

// WithRTBOpenNativeVersion set version
func WithRTBOpenNativeVersion(ver string) BidRequestRTBOption {
	return func(opts *BidRequestRTBOptions) {
		opts.OpenNative.Ver = ver
	}
}

// WithFormatFilter set custom method
func WithFormatFilter(f func(f *types.Format) bool) BidRequestRTBOption {
	return func(opts *BidRequestRTBOptions) {
		opts.FormatFilter = f
	}
}

// WithMaxTimeDuration of the request
func WithMaxTimeDuration(duration time.Duration) BidRequestRTBOption {
	return func(opts *BidRequestRTBOptions) {
		opts.TimeMax = duration
	}
}
