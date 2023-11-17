package security

import (
	"crypto/tls"
)

// NewTLSConfig for security.
func NewTLSConfig(sec Config) (*tls.Config, error) {
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
func NewClientTLSConfig(sec Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(sec.ClientCertFile, sec.ClientKeyFile)
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
