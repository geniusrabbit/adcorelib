package types

import (
	"time"

	"github.com/geniusrabbit/adcorelib/errtype"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/hourstable"
)

// Hours SQL type declaration
type Hours = hourstable.Hours

// Preallocated schedule reject reasons for hot-path Test methods.
var (
	ErrHoursNotAllowed = errtype.Error("hours not allowed")
	ErrDateNotStarted  = errtype.Error("date not started")
	ErrDateEnded       = errtype.Error("date ended")
)

// HoursByString returns hours object by string pattern
func HoursByString(s string) (Hours, error) {
	return hourstable.HoursByString(s)
}

// TimeRequestPointer is interface for request pointer with time
type TimeRequestPointer interface {
	CurrentGeoTime() time.Time
}

// TestTimeRequest checks hours for the request pointer
//
//go:inline
func TestTimeRequest(h Hours, UTCOffset int, pointer TimeRequestPointer) bool {
	// Check current time in ad hours
	if !h.IsAllActive() {
		if UTCOffset >= -12 && UTCOffset <= 14 {
			if !h.TestTime(fasttime.NowUTCPlusOffset(UTCOffset)) {
				return false
			}
		} else if !h.TestTime(pointer.CurrentGeoTime()) {
			return false
		}
	}
	return true
}
