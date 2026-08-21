package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHoursActiveMinutes(t *testing.T) {
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC) // Wednesday noon, minute 0

	t.Run("all active", func(t *testing.T) {
		total, elapsed := HoursActiveMinutes(nil, now)
		assert.Equal(t, 24*60, total)
		assert.Equal(t, 12*60+1, elapsed) // hours 0-11 + current minute
	})

	// 168-bit string: day*24+hour. Wednesday = 3. Hours 0-11 only.
	halfDay := hoursWeekdayPrefix(time.Wednesday, 12)

	t.Run("half day 00-11", func(t *testing.T) {
		total, elapsed := HoursActiveMinutes(halfDay, now)
		assert.Equal(t, 12*60, total)
		assert.Equal(t, 12*60, elapsed) // 00-11 complete; 12:00 is inactive
	})

	t.Run("half day 00-11 at 06:30", func(t *testing.T) {
		at := time.Date(2026, 8, 19, 6, 30, 0, 0, time.UTC)
		total, elapsed := HoursActiveMinutes(halfDay, at)
		assert.Equal(t, 12*60, total)
		assert.Equal(t, 6*60+31, elapsed)
	})
}

func hoursWeekdayPrefix(day time.Weekday, hours int) Hours {
	buf := make([]byte, 7*24)
	for i := range buf {
		buf[i] = '0'
	}
	for hour := 0; hour < hours; hour++ {
		buf[int(day)*24+hour] = '1'
	}
	h, err := HoursByString(string(buf))
	if err != nil {
		panic(err)
	}
	return h
}
