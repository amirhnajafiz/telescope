package config

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config represents the configuration for the application
type Config struct {
	IDP         string `env:"TELESCOPE_BT_IDP" envDefault:"tmp/telescope_idp"`
	IPFSGateway string `env:"TELESCOPE_BT_IPFS_GATEWAY" envDefault:"http://localhost:8080/ipfs/"`
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
