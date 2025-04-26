package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Namespace       = "TELESCOPE"
	SysSubsystem    = "SYS"
	IPFSSubsystem   = "IPFS"
	ClientSubsystem = "CLIENT"
)

// Metrics holds all the Prometheus metrics used in the application
type Metrics struct {
	ClientBandwidth     prometheus.Gauge
	ClientQuality       prometheus.Gauge
	ClientStallRate     prometheus.Counter
	SysErrorCount       *prometheus.CounterVec
	SysBytesTransferred *prometheus.CounterVec
	SysCacheHits        prometheus.Counter
	SysCacheMisses      prometheus.Counter
	SysCacheRatio       prometheus.Gauge
	IPFSBandwidth       prometheus.Gauge
	IPFSRTT             prometheus.Histogram
}

// NewMetrics initializes and returns a new Metrics struct with all the required Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		ClientBandwidth: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "client_bandwidth_bytes",
			Help:      "Client bandwidth usage in bytes.",
			Namespace: Namespace,
			Subsystem: ClientSubsystem,
		}),
		ClientQuality: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "client_quality_index",
			Help:      "Client quality index.",
			Namespace: Namespace,
			Subsystem: ClientSubsystem,
		}),
		ClientStallRate: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "client_stall_rate",
			Help:      "Client stall rate.",
			Namespace: Namespace,
			Subsystem: ClientSubsystem,
		}),
		SysErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "system_error_count",
			Help:      "Number of system errors.",
			Namespace: Namespace,
			Subsystem: SysSubsystem,
		}, []string{"path"}),
		SysBytesTransferred: promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "system_bytes_transferred",
			Help:      "Total bytes transferred by the system.",
			Namespace: Namespace,
			Subsystem: SysSubsystem,
		}, []string{"path"}),
		SysCacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "system_cache_hits",
			Help:      "Number of cache hits.",
			Namespace: Namespace,
			Subsystem: SysSubsystem,
		}),
		SysCacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Name:      "system_cache_misses",
			Help:      "Number of cache misses.",
			Namespace: Namespace,
			Subsystem: SysSubsystem,
		}),
		SysCacheRatio: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "system_cache_ratio",
			Help:      "Cache hit ratio.",
			Namespace: Namespace,
			Subsystem: SysSubsystem,
		}),
		IPFSBandwidth: promauto.NewGauge(prometheus.GaugeOpts{
			Name:      "ipfs_bandwidth_bytes",
			Help:      "IPFS bandwidth usage in bytes.",
			Namespace: Namespace,
			Subsystem: IPFSSubsystem,
		}),
		IPFSRTT: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:      "round_trip_time",
			Help:      "Round trip time.",
			Namespace: Namespace,
			Subsystem: IPFSSubsystem,
		}),
	}
}
