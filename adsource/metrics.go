package adsource

import (
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"

	"github.com/demdxx/gocast/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics wrapper
type Metrics struct {
	requestCounter *prometheus.CounterVec
}

// IncrementBidRequestCount metric
func (m *Metrics) IncrementBidRequestCount(source adtype.Source, request adtype.BidRequester, duration time.Duration) {
	m.counter().WithLabelValues(
		gocast.Str(source.ID()),
		gocast.Str(request.TargetID()),
		gocast.Str(request.DeviceInfo().DeviceType),
		request.AuctionType().Name(),
		"0",
	).Inc()
}

// IncrementBidErrorCount metric
func (m *Metrics) IncrementBidErrorCount(source adtype.Source, request adtype.BidRequester, err error) {
	m.counter().WithLabelValues(
		gocast.Str(source.ID()),
		gocast.Str(request.TargetID()),
		gocast.Str(request.DeviceInfo().DeviceType),
		request.AuctionType().Name(),
		gocast.IfThen(err != nil, "1", "0"),
	).Inc()
}

func (m *Metrics) counter() *prometheus.CounterVec {
	if m.requestCounter == nil {
		m.requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "multisource_count",
			Help: "Count of requests by source params",
		}, []string{"source_id", "zone_id", "device_type", "auction_type", "error"})
	}
	return m.requestCounter
}
