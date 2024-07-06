package loader

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/fasttime"
)

// LoaderFnk type
type LoaderFnk func(objectTarget any, lastUpdate *time.Time) error

// LoaderTarget object interface
type LoaderTarget interface {
	Result() []any
	Reset()
}

// DataAccessor with data loading
type DataAccessor interface {
	Data() ([]any, error)
	NeedUpdate() bool
}

// PeriodicDataAccessor with interval between reloeads
type PeriodicDataAccessor struct {
	mx         sync.Mutex
	lastUpdate uint64
	fullReload uint64
	period     uint64
	target     LoaderTarget
	loader     LoaderFnk
	data       []any

	metricReloadCounter *prometheus.CounterVec
	metricLoadedCount   prometheus.Gauge
}

// NewPeriodicFullreloader accessor
func NewPeriodicFullreloader(target LoaderTarget, loader LoaderFnk, period time.Duration, metric string) *PeriodicDataAccessor {
	return NewPeriodicReloader(target, loader, period, 0, metric)
}

// NewPeriodicReloader accessor
func NewPeriodicReloader(target LoaderTarget, loader LoaderFnk, period, fulleReloadPeriod time.Duration, metric string) *PeriodicDataAccessor {
	if period.Seconds() < 1 {
		panic("time period less the 1s")
	}
	return &PeriodicDataAccessor{
		lastUpdate: 0,
		fullReload: uint64(fulleReloadPeriod.Seconds()),
		period:     uint64(period.Seconds()),
		target:     target,
		loader:     loader,
		data:       nil,

		metricReloadCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: metric + "_count",
			Help: "Count reloads",
		}, []string{"full"}),
		metricLoadedCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: metric + "_number",
			Help: "Namber value",
		}),
	}
}

func (l *PeriodicDataAccessor) NeedUpdate() bool {
	return atomic.LoadUint64(&l.lastUpdate) < fasttime.UnixTimestamp()-l.period
}

// Data returns loaded data and reload if necessary
func (l *PeriodicDataAccessor) Data() ([]any, error) {
	if l.NeedUpdate() {
		l.mx.Lock()
		defer l.mx.Unlock()
		if l.NeedUpdate() {
			lastUpdate := time.Unix(int64(l.lastUpdate), 0)
			if l.fullReload == 0 || uint64(fasttime.Now().Sub(lastUpdate).Seconds()) > l.fullReload {
				l.target.Reset()
				l.metricReloadCounter.WithLabelValues("1").Inc()
				lastUpdate = lastUpdate.AddDate(-10, 0, 0)
			} else {
				l.metricReloadCounter.WithLabelValues("0").Inc()
			}
			if err := l.loader(l.target, &lastUpdate); err != nil {
				return nil, err
			}
			l.data = l.target.Result()
			l.metricLoadedCount.Set(float64(len(l.data)))
			zap.L().Debug("data loaded",
				zap.String("model", reflect.TypeOf(l.target).String()),
				zap.Int("count", len(l.data)),
			)
			atomic.StoreUint64(&l.lastUpdate, fasttime.UnixTimestamp())
		}
	}
	return l.data, nil
}
