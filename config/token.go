package config

import (
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
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

func tokenSSHConfig(cfg *Config) *ssh.Config {
	if !token.IsEnabled(cfg.Token) {
		return nil
	}

	return cfg.Token.SSH
}
