package counter

import (
	"math"
	"math/rand"
	"sync/atomic"
)

const treshold = 1000.0

// ErrorCounter of errors
type ErrorCounter struct {
	value int32
}

// Inc - rement counter value
func (cn *ErrorCounter) Inc() {
	if newVal := atomic.AddInt32(&cn.value, 1); newVal > treshold {
		atomic.StoreInt32(&cn.value, treshold)
	}
}

// Dec - rement counter value
func (cn *ErrorCounter) Dec() {
	if newVal := atomic.AddInt32(&cn.value, -1); newVal < -treshold {
		atomic.StoreInt32(&cn.value, -treshold)
	}
}

// Do counter incrementation
func (cn *ErrorCounter) Do(inc bool) {
	if inc {
		cn.Inc()
	} else {
		cn.Dec()
	}
}

// Next value
func (cn *ErrorCounter) Next() bool {
	return rand.Float64()*1.03 > cn.factor()
}

func (cn *ErrorCounter) factor() float64 {
	return normalise(cn.getValue() / treshold * 7.8)
}

func (cn *ErrorCounter) getValue() float64 {
	return float64(atomic.LoadInt32(&cn.value))
}

func normalise(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(x - 10)))
}
