package config

import (
	"fmt"
	"os"
)

const (
	httpHostEnvName = "HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func (ac *httpConfig) Address() string {
	return fmt.Sprintf("%s:%s", "", ac.port)
}

func NewHttpConfig() (Config, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("No host name in .env")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("No port in .env")
	}

	config := httpConfig{
		host: host,
		port: port,
	}

	return &config, nil
}
