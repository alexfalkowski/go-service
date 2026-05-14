package config

import (
	"crypto/x509"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewCertPool resolves and parses cfg's configured CA bundle.
func NewCertPool(fs *os.FS, cfg *Config) (*x509.CertPool, error) {
	ca, err := cfg.GetCA(fs)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return nil, ErrInvalidCA
	}

	return pool, nil
}
