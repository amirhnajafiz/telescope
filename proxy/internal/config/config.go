package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config represents the configuration for the application
type Config struct {
	Port        int    `env:"TELESCOPE_PORT" envDefault:"8080"`
	MetricsPort int    `env:"TELESCOPE_METRICS_PORT" envDefault:"9090"`
	Debug       bool   `env:"TELESCOPE_DEBUG" envDefault:"false"`
	CachePath   string `env:"TELESCOPE_CACHE_PATH" envDefault:"tmp/telescope-cache"`
	Jaeger      string `env:"TELESCOPE_JAEGER" envDefault:""`
	IPFSGateway string `env:"TELESCOPE_IPFS_GATEWAY" envDefault:"http://localhost:5001/"`
}

// LoadConfigs loads the configuration from environment variables
func LoadConfigs() (*Config, error) {
	// check if ".env" file exists
	if _, err := os.Stat(".env"); err == nil {
		// load the ".env" file
		if err := godotenv.Load(".env"); err != nil {
			return nil, err
		}
	}

	// load the configuration from environment variables
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
