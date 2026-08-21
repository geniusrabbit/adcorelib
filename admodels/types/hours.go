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
	ErrHoursNotAllowed       = errtype.Error("hours not allowed")
	ErrDateNotStarted        = errtype.Error("date not started")
	ErrDateEnded             = errtype.Error("date ended")
	ErrBudgetValueNotAllowed = errtype.Error("budget value not allowed")
	ErrPacingNotAllowed      = errtype.Error("pacing not allowed")
)

// HoursByString returns hours object by string pattern
func HoursByString(s string) (Hours, error) {
	return hourstable.HoursByString(s)
}

// TimeRequestPointer is interface for request pointer with time
type TimeRequestPointer interface {
	CurrentGeoTime() time.Time
}

// TimeRequestNow returns the clock used for hours / EVENLY pacing.
// Fixed UTCOffset in [-12, 14] uses NowUTCPlusOffset; otherwise geo time.
func TimeRequestNow(UTCOffset int, pointer TimeRequestPointer) time.Time {
	if UTCOffset >= -12 && UTCOffset <= 14 {
		return fasttime.NowUTCPlusOffset(UTCOffset)
	}
	if pointer != nil {
		return pointer.CurrentGeoTime()
	}
	return fasttime.NowUTC()
}

// HoursActiveMinutes returns total active minutes for t's weekday and elapsed
// active minutes from midnight up to and including the current minute.
// Empty / all-active hours count as 24h (1440 minutes).
func HoursActiveMinutes(h Hours, t time.Time) (total, elapsed int) {
	weekday := t.Weekday()
	nowHour, nowMin := t.Hour(), t.Minute()
	allActive := h.IsAllActive()
	for hour := 0; hour < 24; hour++ {
		if !allActive && !h.TestHour(weekday, byte(hour)) {
			continue
		}
		total += 60
		if hour < nowHour {
			elapsed += 60
		} else if hour == nowHour {
			elapsed += nowMin + 1
		}
	}
	return total, elapsed
}

// TestTimeRequest checks hours for the request pointer
//
//go:inline
func TestTimeRequest(h Hours, UTCOffset int, pointer TimeRequestPointer) bool {
	// Check current time in ad hours
	if !h.IsAllActive() {
		if !h.TestTime(TimeRequestNow(UTCOffset, pointer)) {
			return false
		}
	}
	return true
}
