package config

import (
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
)

func cryptoAESConfig(cfg *Config) *aes.Config {
	if !cfg.Crypto.IsEnabled() {
		return nil
	}
	return cfg.Crypto.AES
}

func cryptoED25519Config(cfg *Config) *ed25519.Config {
	if !cfg.Crypto.IsEnabled() {
		return nil
	}
	return cfg.Crypto.Ed25519
}

func cryptoHMACConfig(cfg *Config) *hmac.Config {
	if !cfg.Crypto.IsEnabled() {
		return nil
	}
	return cfg.Crypto.HMAC
}

func cryptoRSAConfig(cfg *Config) *rsa.Config {
	if !cfg.Crypto.IsEnabled() {
		return nil
	}
	return cfg.Crypto.RSA
}

func cryptoSSHConfig(cfg *Config) *ssh.Config {
	if !cfg.Crypto.IsEnabled() {
		return nil
	}
	return cfg.Crypto.SSH
}
