package config

import (
	"fmt"
	"os"
)

const (
	grpcHostEnvName = "HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func (ac *grpcConfig) Address() string {
	return fmt.Sprintf("%s:%s", "", ac.port)
}

func NewGrpcConfig() (Config, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("No host name in .env")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("No port in .env")
	}

	config := grpcConfig{
		host: host,
		port: port,
	}

	return &config, nil
}
