package security

import (
	"crypto/tls"
)

// NewTLSConfig for security.
func NewTLSConfig(sec *Config) (*tls.Config, error) {
	c := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if sec == nil {
		return c, nil
	}

	cert, err := tls.LoadX509KeyPair(sec.CertFile, sec.KeyFile)
	if err != nil {
		return nil, err
	}

	c.Certificates = []tls.Certificate{cert}

	return c, nil
}

// NewClientTLSConfig for security.
func NewClientTLSConfig(sec *Config) (*tls.Config, error) {
	c := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if sec == nil {
		return c, nil
	}

	cert, err := tls.LoadX509KeyPair(sec.CertFile, sec.KeyFile)
	if err != nil {
		return nil, err
	}

	c.Certificates = []tls.Certificate{cert}

	return c, nil
}
