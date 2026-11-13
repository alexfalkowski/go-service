package config

import (
	"github.com/alexfalkowski/go-service/v3/crypto/aes"
	"github.com/alexfalkowski/go-service/v3/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v3/crypto/hmac"
	"github.com/alexfalkowski/go-service/v3/crypto/rsa"
)

func cryptoAESConfig(cfg *Config) *aes.Config {
	if cfg.Crypto.IsEnabled() {
		return cfg.Crypto.AES
	}
	return nil
}

func cryptoED25519Config(cfg *Config) *ed25519.Config {
	if cfg.Crypto.IsEnabled() {
		return cfg.Crypto.Ed25519
	}
	return nil
}

func cryptoHMACConfig(cfg *Config) *hmac.Config {
	if cfg.Crypto.IsEnabled() {
		return cfg.Crypto.HMAC
	}
	return nil
}

func cryptoRSAConfig(cfg *Config) *rsa.Config {
	if cfg.Crypto.IsEnabled() {
		return cfg.Crypto.RSA
	}
	return nil
}
