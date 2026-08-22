package config

import (
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

func grpcConfig(cfg *Config) *grpc.Config {
	if cfg.Transport.IsEnabled() {
		return cfg.Transport.GRPC
	}
	return nil
}

func httpConfig(cfg *Config) *http.Config {
	if cfg.Transport.IsEnabled() {
		return cfg.Transport.HTTP
	}
	return nil
}

func accessConfig(cfg *Config) *access.Config {
	if cfg.Transport.IsEnabled() {
		return cfg.Transport.Access
	}
	return nil
}
