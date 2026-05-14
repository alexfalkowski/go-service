package config

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewKeyPair resolves and parses cfg's configured leaf certificate/private-key pair.
func NewKeyPair(fs *os.FS, cfg *Config) (tls.Certificate, error) {
	cert, err := cfg.GetCert(fs)
	if err != nil {
		return tls.Certificate{}, err
	}

	key, err := cfg.GetKey(fs)
	if err != nil {
		return tls.Certificate{}, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return pair, nil
}
