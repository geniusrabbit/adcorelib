package fasttime

import (
	"sync/atomic"
	"time"
)

var currentTimestamp = uint64(time.Now().UnixNano())

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

// Now return  current time object
func Now() time.Time {
	return time.Unix(0, int64(UnixTimestampNano()))
}

// NowUTC return  current time object in UTC
func NowUTC() time.Time {
	return time.Unix(0, int64(UnixTimestampNano())).UTC()
}

// NowUTCPlusOffset return  current time object in UTC with offset in hours
func NowUTCPlusOffset(offset int) time.Time {
	tm := NowUTC()
	if offset == 0 {
		return tm
	}
	return tm.Add(time.Duration(offset) * time.Hour)
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
