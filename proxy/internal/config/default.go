package config

// Default returns the default configuration for the application
func Default() *Config {
	return &Config{
		Port:        3000,
		Debug:       true,
		Jaeger:      "",
		MetricsPort: 0,
		CachePath:   "tmp/telescope-cache",
		IPFSGateway: "http://localhost:8080/ipfs/",
	}
}
