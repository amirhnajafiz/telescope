package config

// Default returns the default configuration for the application
func Default() *Config {
	return &Config{
		IDP:         "tmp/telescope_idp",
		IPFSGateway: "http://localhost:8080/ipfs/",
	}
}
