package config

import (
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

func tokenConfig(cfg *Config) *token.Config {
	return cfg.Token
}

func tokenAccessConfig(cfg *Config) *access.Config {
	if !cfg.Token.IsEnabled() {
		return nil
	}
	return cfg.Token.Access
}

func tokenJWTConfig(cfg *Config) *jwt.Config {
	if !cfg.Token.IsEnabled() {
		return nil
	}
	return cfg.Token.JWT
}

func tokenPasetoConfig(cfg *Config) *paseto.Config {
	if !cfg.Token.IsEnabled() {
		return nil
	}
	return cfg.Token.Paseto
}

func tokenSSHConfig(cfg *Config) *ssh.Config {
	if !cfg.Token.IsEnabled() {
		return nil
	}
	return cfg.Token.SSH
}
