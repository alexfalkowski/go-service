package security

import (
	"crypto/tls"
	"encoding/base64"
)

// NewTLSConfig for security.
func NewTLSConfig(sec *Config) (*tls.Config, error) {
	c := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !IsEnabled(sec) || !sec.HasKeyPair() {
		return c, nil
	}

	dc, err := base64.StdEncoding.DecodeString(sec.GetCert())
	if err != nil {
		return c, err
	}

	dk, err := base64.StdEncoding.DecodeString(sec.GetKey())
	if err != nil {
		return c, err
	}

	cert, err := tls.X509KeyPair(dc, dk)
	if err != nil {
		return c, err
	}

	c.Certificates = []tls.Certificate{cert}

	return c, nil
}
