package security

import (
	"crypto/tls"
)

// NewTLSConfig for security.
//
//nolint:nilnil
func NewTLSConfig(sec *Config) (*tls.Config, error) {
	if sec == nil {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(sec.CertFile, sec.KeyFile)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return conf, nil
}

// NewClientTLSConfig for security.
//
//nolint:nilnil
func NewClientTLSConfig(sec *Config) (*tls.Config, error) {
	if sec == nil {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(sec.CertFile, sec.KeyFile)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return conf, nil
}
