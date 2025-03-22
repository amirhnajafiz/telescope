package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics used in the application
type Metrics struct {
	RoundTripTime     *prometheus.HistogramVec
	RequestDuration   *prometheus.HistogramVec
	Bandwidth         *prometheus.CounterVec
	ErrorCount        *prometheus.CounterVec
	BytesTransferred  *prometheus.CounterVec
	NumberOfStreams   prometheus.Counter
	CacheHits         prometheus.Counter
	CacheMisses       prometheus.Counter
	CacheRatio        prometheus.Gauge
	LocalStorageSize  prometheus.Gauge
	ActiveConnections prometheus.Gauge
}

// NewMetrics initializes and returns a new Metrics struct with all the required Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		RoundTripTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "round_trip_time_seconds",
			Help: "Round trip time in seconds.",
		}, []string{"method", "endpoint"}),
		Bandwidth: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "bandwidth_bytes",
			Help: "Bandwidth usage in bytes.",
		}, []string{"method", "endpoint"}),
		NumberOfStreams: promauto.NewCounter(prometheus.CounterOpts{
			Name: "number_of_streams",
			Help: "Number of streams.",
		}),
		CacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits",
			Help: "Number of cache hits.",
		}),
		CacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses",
			Help: "Number of cache misses.",
		}),
		CacheRatio: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cache_ratio",
			Help: "Cache hit ratio.",
		}),
		LocalStorageSize: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "local_storage_size_bytes",
			Help: "Local storage size in bytes.",
		}),
		RequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Duration of requests in seconds.",
		}, []string{"method", "endpoint"}),
		ErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "error_count",
			Help: "Number of errors.",
		}, []string{"method", "endpoint"}),
		ActiveConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections.",
		}),
		BytesTransferred: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "bytes_transferred",
			Help: "Total bytes transferred.",
		}, []string{"method", "endpoint"}),
	}
}
