//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package auction

import (
	"sync"

	"geniusrabbit.dev/adcorelib/adtype"
)

var (
	countersPool = sync.Pool{
		New: func() any { return &counters{items: make([]int, 0, 10)} },
	}
	responseListPool = sync.Pool{
		New: func() any { return make([]adtype.ResponserItemCommon, 0, 10) },
	}
)

// borrowCounters object
func borrowCounters() *counters {
	return countersPool.Get().(*counters)
}

// returnCounter back to pool
func returnCounter(arr *counters) {
	if arr != nil {
		if len(arr.items) > 0 {
			arr.items = arr.items[:0]
		}
		countersPool.Put(arr)
	}
}

// borrowResponseList object
func borrowResponseList() []adtype.ResponserItemCommon {
	return responseListPool.Get().([]adtype.ResponserItemCommon)
}

// returnResponseList back to pool
func returnResponseList(arr []adtype.ResponserItemCommon) {
	if arr != nil {
		responseListPool.Put(arr[:0])
	}
}

type counters struct {
	items []int
}

func (c counters) count(idx int) int {
	if idx >= len(c.items) {
		return 0
	}
	return c.items[idx]
}

func (c *counters) inc(idx, v int) *counters {
	if v != 0 {
		for idx >= len(c.items) {
			c.items = append(c.items, 0)
		}
		c.items[idx] += v
	}
	return c
}
