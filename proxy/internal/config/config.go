package config

import "github.com/caarlos0/env/v10"

// Config represents the configuration for the application
type Config struct {
	Port        int    `env:"TELESCOPE_PORT" envDefault:"8080"`
	Debug       bool   `env:"TELESCOPE_DEBUG" envDefault:"false"`
	Jaeger      string `env:"TELESCOPE_JAEGER" envDefault:""`
	MetricsPort int    `env:"TELESCOPE_METRICS_PORT" envDefault:"9090"`
	IPFSGateway string `env:"TELESCOPE_IPFS_GATEWAY" envDefault:"http://localhost:8080/ipfs/"`
}

// LoadConfigs loads the configuration from environment variables
func LoadConfigs() (*Config, error) {
	// load the configuration from environment variables
	cfg := Default()
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
