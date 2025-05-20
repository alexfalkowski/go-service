package config

import (
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	cs "github.com/alexfalkowski/go-service/v2/crypto/ssh"
)

func cryptoAESConfig(cfg *Config) *aes.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.AES
}

func cryptoED25519Config(cfg *Config) *ed25519.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.Ed25519
}

func cryptoHMACConfig(cfg *Config) *hmac.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.HMAC
}

func cryptoRSAConfig(cfg *Config) *rsa.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.RSA
}

func cryptoSSHConfig(cfg *Config) *cs.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.SSH
}
