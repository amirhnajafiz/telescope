package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace = "telescope"
	Subsystem = "proxy"
)

// Metrics holds all the Prometheus metrics used in the application
type Metrics struct {
	RoundTripTime    *prometheus.HistogramVec
	Bandwidth        *prometheus.CounterVec
	ErrorCount       *prometheus.CounterVec
	BytesTransferred *prometheus.CounterVec
	CacheHits        prometheus.Counter
	CacheMisses      prometheus.Counter
	CacheRatio       prometheus.Gauge
	LocalStorageSize prometheus.Gauge
}

// NewMetrics initializes and returns a new Metrics struct with all the required Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		RoundTripTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:      "round_trip_time_seconds",
			Help:      "Round trip time in seconds.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}, []string{"method", "endpoint"}),
		Bandwidth: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "bandwidth_bytes",
			Help:      "Bandwidth usage in bytes.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}, []string{"method", "endpoint"}),
		CacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "cache_hits",
			Help:      "Number of cache hits.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}),
		CacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "cache_misses",
			Help:      "Number of cache misses.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}),
		CacheRatio: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "cache_ratio",
			Help:      "Cache hit ratio.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}),
		LocalStorageSize: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "local_storage_size_bytes",
			Help:      "Local storage size in bytes.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}),
		ErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "error_count",
			Help:      "Number of errors.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}, []string{"method", "endpoint"}),
		BytesTransferred: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "bytes_transferred",
			Help:      "Total bytes transferred.",
			Namespace: Namespace,
			Subsystem: Subsystem,
		}, []string{"method", "endpoint"}),
	}
}
