package prometheuswrapper

import (
	"time"

	"geniusrabbit.dev/corelib/openlatency"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Wrapper struct {
	counter  *openlatency.MetricsCounter
	queries  prometheus.Counter
	errors   prometheus.Counter
	noBid    prometheus.Counter
	skip     prometheus.Counter
	timeouts prometheus.Counter
	success  prometheus.Counter
	latency  prometheus.Observer
}

func NewWrapperDefault(prefix string, tags, params []string) *Wrapper {
	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	return &Wrapper{
		counter: openlatency.NewMetricsCounter(),
		queries: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "queries_count",
			Help: "Count of requests",
		}, tags).WithLabelValues(params...),
		errors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "errors_count",
			Help: "Count of errors",
		}, tags).WithLabelValues(params...),
		noBid: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "nobid_count",
			Help: "Count of nobids",
		}, tags).WithLabelValues(params...),
		skip: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "skip_count",
			Help: "Count of skips",
		}, tags).WithLabelValues(params...),
		timeouts: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "timeout_count",
			Help: "Count of timeouts",
		}, tags).WithLabelValues(params...),
		success: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "success_count",
			Help: "Count of success",
		}, tags).WithLabelValues(params...),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "latency_seconds",
			Help:    "Histogram of response time in seconds",
			Buckets: buckets,
		}, tags).WithLabelValues(params...),
	}
}

// UpdateQueryLatency of request
func (wrp *Wrapper) UpdateQueryLatency(latency time.Duration) {
	wrp.latency.Observe(latency.Seconds())
	wrp.counter.UpdateQueryLatency(latency)
}

// BeginQuery new query counter
func (wrp *Wrapper) BeginQuery() int32 {
	wrp.queries.Inc()
	return wrp.counter.BeginQuery()
}

// IncTimeout counter
func (wrp *Wrapper) IncTimeout() int32 {
	wrp.timeouts.Inc()
	return wrp.counter.IncTimeout()
}

// IncNobid counter
func (wrp *Wrapper) IncNobid() int32 {
	wrp.noBid.Inc()
	return wrp.counter.IncNobid()
}

// IncSkip counter
func (wrp *Wrapper) IncSkip() int32 {
	wrp.skip.Inc()
	return wrp.counter.IncSkip()
}

// IncSuccess counter
func (wrp *Wrapper) IncSuccess() int32 {
	wrp.success.Inc()
	return wrp.counter.IncSuccess()
}

// IncError counter
func (wrp *Wrapper) IncError(etype openlatency.MetricErrorType, code string) {
	wrp.errors.Inc()
	wrp.counter.IncError(etype, code)
}

// FillMetrics info object
func (wrp *Wrapper) FillMetrics(info *openlatency.MetricsInfo) {
	wrp.counter.FillMetrics(info)
}
