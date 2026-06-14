package config

import (
	"crypto/x509"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewCertPool resolves and parses cfg's configured CA bundle.
//
// It returns [ErrInvalidCA] when cfg is nil, CA is not configured, or the resolved CA bytes cannot be
// appended to a certificate pool. Source-resolution errors are returned directly.
func NewCertPool(fs *os.FS, cfg *Config) (*x509.CertPool, error) {
	if !cfg.HasCA() {
		return nil, ErrInvalidCA
	}

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
