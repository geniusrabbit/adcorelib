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
func (m *Metrics) IncrementBidRequestCount(source adtype.Source, request *adtype.BidRequest, duration time.Duration) {
	m.counter().WithLabelValues(
		gocast.ToString(source.ID()),
		gocast.ToString(request.TargetID()),
		gocast.ToString(request.DeviceInfo().DeviceType),
		request.AuctionType.Name(),
		"0",
	).Inc()
}

// IncrementBidErrorCount metric
func (m *Metrics) IncrementBidErrorCount(source adtype.Source, request *adtype.BidRequest, err error) {
	m.counter().WithLabelValues(
		gocast.ToString(source.ID()),
		gocast.ToString(request.TargetID()),
		gocast.ToString(request.DeviceInfo().DeviceType),
		request.AuctionType.Name(),
		isToBinS(err != nil),
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

func isToBinS(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
