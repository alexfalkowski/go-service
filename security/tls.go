package security

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

// ErrBadCertificate for security.
var ErrBadCertificate = errors.New("bad certificate")

// NewServerTLSConfig for security.
func NewServerTLSConfig(sec *Config) (*tls.Config, error) {
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

// NewClientTLSConfig for security.
func NewClientTLSConfig(sec *Config) (*tls.Config, error) {
	c := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !IsEnabled(sec) || !sec.HasCert() {
		return c, nil
	}

	dc, err := base64.StdEncoding.DecodeString(sec.GetCert())
	if err != nil {
		return c, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(dc) {
		return c, ErrBadCertificate
	}

	c.RootCAs = cp

	return c, nil
}
