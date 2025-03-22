package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metrics holds all the Prometheus metrics used in the application
type Metrics struct {
	RoundTripTime     *prometheus.HistogramVec
	Bandwidth         *prometheus.CounterVec
	NumberOfStreams   prometheus.Counter
	CacheHits         prometheus.Counter
	CacheMisses       prometheus.Counter
	CacheRatio        prometheus.Gauge
	LocalStorageSize  prometheus.Gauge
	RequestDuration   *prometheus.HistogramVec
	ErrorCount        *prometheus.CounterVec
	ActiveConnections prometheus.Gauge
	BytesTransferred  *prometheus.CounterVec
	CPUUsage          prometheus.Gauge
	MemoryUsage       prometheus.Gauge
	DiskIO            *prometheus.CounterVec
	NetworkIO         *prometheus.CounterVec
}

// NewMetrics initializes and returns a new Metrics struct with all the required Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		RoundTripTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "round_trip_time_seconds",
			Help: "Round trip time in seconds.",
		}, []string{"method", "endpoint"}),
		Bandwidth: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "bandwidth_bytes",
			Help: "Bandwidth usage in bytes.",
		}, []string{"method", "endpoint"}),
		NumberOfStreams: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "number_of_streams",
			Help: "Number of streams.",
		}),
		CacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits",
			Help: "Number of cache hits.",
		}),
		CacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses",
			Help: "Number of cache misses.",
		}),
		CacheRatio: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cache_ratio",
			Help: "Cache hit ratio.",
		}),
		LocalStorageSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "local_storage_size_bytes",
			Help: "Local storage size in bytes.",
		}),
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Duration of requests in seconds.",
		}, []string{"method", "endpoint"}),
		ErrorCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "error_count",
			Help: "Number of errors.",
		}, []string{"method", "endpoint"}),
		ActiveConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections.",
		}),
		BytesTransferred: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "bytes_transferred",
			Help: "Total bytes transferred.",
		}, []string{"method", "endpoint"}),
		CPUUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "CPU usage.",
		}),
		MemoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage.",
		}),
		DiskIO: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "disk_io_operations",
			Help: "Disk I/O operations.",
		}, []string{"operation"}),
		NetworkIO: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "network_io_bytes",
			Help: "Network I/O in bytes.",
		}, []string{"direction"}),
	}
}
