package fasttime

import (
	"sync/atomic"
	"time"
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Second / 10)
		defer ticker.Stop()
		for tm := range ticker.C {
			t := uint64(tm.UnixNano())
			atomic.StoreUint64(&currentTimestamp, t)
		}
	}()
}

var currentTimestamp = uint64(time.Now().UnixNano())

// Now return  current time object
func Now() time.Time {
	return time.Unix(0, int64(UnixTimestampNano()))
}

// UnixTimestampNano returns the current unix timestamp in nanoseconds.
//
// It is faster than time.Now().UnixNano()
func UnixTimestampNano() uint64 {
	return atomic.LoadUint64(&currentTimestamp)
}

// UnixTimestamp returns the current unix timestamp in seconds.
//
// It is faster than time.Now().Unix()
func UnixTimestamp() uint64 {
	return UnixTimestampNano() / uint64(time.Second)
}

// UnixDate returns date from the current unix timestamp.
//
// The date is calculated by dividing unix timestamp by (24*3600)
func UnixDate() uint64 {
	return UnixTimestamp() / (24 * 3600)
}

// UnixHour returns hour from the current unix timestamp.
//
// The hour is calculated by dividing unix timestamp by 3600
func UnixHour() uint64 {
	return UnixTimestamp() / 3600
}
