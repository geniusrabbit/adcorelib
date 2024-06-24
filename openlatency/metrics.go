package openlatency

// MetricErrorType values
type MetricErrorType string

// Error type list...
const (
	MetricErrorHTTP    MetricErrorType = "http"
	MetricErrorNetwork MetricErrorType = "network"
)

type MetricErrorRate struct {
	Type MetricErrorType `json:"type"`
	Code string          `json:"code"`
	Rate float64         `json:"rate"`
}

type MetricsGeoRate struct {
	Country string  `json:"country"`
	Rate    float64 `json:"rate"`
}

// MetricsInfo describes basic metric information of AdNetworks integration
// All counters it's numbers per second
type MetricsInfo struct {
	ID         uint64            `json:"id"`
	Protocol   string            `json:"protocol"`
	Codename   string            `json:"codename,omitempty"`
	Traceroute string            `json:"traceroute,omitempty"`
	MinLatency int64             `json:"min_latency_ms"` // Minimal request delay in Millisecond
	MaxLatency int64             `json:"max_latency_ms"` // Maximal request delay in Millisecond
	AvgLatency int64             `json:"avg_latency_ms"` // Average request delay in Millisecond
	QPSLimit   int               `json:"qps_limit,omitempty"`
	QPS        float64           `json:"qps"`
	Skips      float64           `json:"skips_qps"`
	Success    float64           `json:"success_qps"`
	Timeouts   float64           `json:"timeouts_qps"`
	NoBids     float64           `json:"no_bids_qps"`
	Errors     float64           `json:"errors_qps"`
	ErrorRates []MetricErrorRate `json:"error_rates,omitempty"`
	GeoRates   []MetricsGeoRate  `json:"geo_rates,omitempty"`
}
