package tls

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewConfig constructs a `crypto/tls`.Config from cfg.
//
// It always sets a minimum TLS version of 1.2 and requires/validates client
// certificates (mTLS) via `tls.RequireAndVerifyClientCert`.
//
// If cfg is nil or does not contain both a certificate and key, a config with
// those defaults is returned without loading any key material.
//
// Cert and key bytes are read via the provided filesystem using the configured
// "source strings" (literal, `file:`, or `env:`).
func NewConfig(fs *os.FS, cfg *Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !cfg.IsEnabled() || !cfg.HasKeyPair() {
		return config, nil
	}

	cert, err := cfg.GetCert(fs)
	if err != nil {
		return config, err
	}

	key, err := cfg.GetKey(fs)
	if err != nil {
		return config, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return config, err
	}

	config.Certificates = []tls.Certificate{pair}
	return config, nil
}
