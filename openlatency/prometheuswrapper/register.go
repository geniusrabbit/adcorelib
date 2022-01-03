package prometheuswrapper

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	regMx     sync.Mutex
	counters  = map[string]*prometheus.CounterVec{}
	histogram = map[string]*prometheus.HistogramVec{}
)

func newCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	regMx.Lock()
	defer regMx.Unlock()
	cnt := counters[opts.Name]
	if cnt != nil {
		return cnt
	}
	cnt = promauto.NewCounterVec(opts, labelNames)
	counters[opts.Name] = cnt
	return cnt
}

func newHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	regMx.Lock()
	defer regMx.Unlock()
	hist := histogram[opts.Name]
	if hist != nil {
		return hist
	}
	hist = prometheus.NewHistogramVec(opts, labelNames)
	histogram[opts.Name] = hist
	return hist
}
