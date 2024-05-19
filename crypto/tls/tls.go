package tls

import (
	"crypto/tls"
)

// NewConfig for tls.
func NewConfig(cfg *Config) (*tls.Config, error) {
	c := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !IsEnabled(cfg) || !cfg.HasKeyPair() {
		return c, nil
	}

	dc, err := cfg.GetCert()
	if err != nil {
		return c, err
	}

	dk, err := cfg.GetKey()
	if err != nil {
		return c, err
	}

	pair, err := tls.X509KeyPair(dc, dk)
	if err != nil {
		return c, err
	}

	c.Certificates = []tls.Certificate{pair}

	return c, nil
}
