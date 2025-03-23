package config

import (
	"fmt"
	"os"
)

const (
	authHostEnvName = "AUTH_HOST"
	authPortEnvName = "AUTH_PORT"
)

type AuthConfig interface {
	Address() string
}

type authConfig struct {
	host string
	port string
}

func (ac *authConfig) Address() string {
	return fmt.Sprintf(":%s", ac.port)
}

func NewAuthConfig() (AuthConfig, error) {
	host := os.Getenv(authHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("No host name in .env")
	}

	port := os.Getenv(authPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("No port in .env")
	}

	config := authConfig{
		host: host,
		port: port,
	}

	return &config, nil
}
