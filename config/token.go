package config

import (
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/token/jwt"
	"github.com/alexfalkowski/go-service/token/paseto"
	ts "github.com/alexfalkowski/go-service/token/ssh"
)

func tokenConfig(cfg *Config) *token.Config {
	if !token.IsEnabled(cfg.Token) {
		return nil
	}

	return cfg.Token
}

func tokenJWTConfig(cfg *Config) *jwt.Config {
	if !token.IsEnabled(cfg.Token) {
		return nil
	}

	return cfg.Token.JWT
}

func tokenPasetoConfig(cfg *Config) *paseto.Config {
	if !token.IsEnabled(cfg.Token) {
		return nil
	}

	return cfg.Token.Paseto
}

func tokenSSHConfig(cfg *Config) *ts.Config {
	if !token.IsEnabled(cfg.Token) {
		return nil
	}

	return cfg.Token.SSH
}
