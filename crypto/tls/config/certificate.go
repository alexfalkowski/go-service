package config

import (
	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewKeyPair resolves and parses cfg's configured leaf certificate/private-key pair.
//
// It returns [github.com/alexfalkowski/go-service/v2/crypto/errors.ErrMissingKey] when cfg is nil, either source
// string is missing, or either source resolves to empty bytes. Otherwise it returns source-resolution
// errors or [github.com/alexfalkowski/go-service/v2/crypto/tls.X509KeyPair] parsing errors.
func NewKeyPair(fs *os.FS, cfg *Config) (tls.Certificate, error) {
	if cfg == nil || strings.IsEmpty(cfg.Cert) || strings.IsEmpty(cfg.Key) {
		return tls.Certificate{}, crypto.ErrMissingKey
	}

	cert, err := cfg.GetCert(fs)
	if err != nil {
		return tls.Certificate{}, err
	}
	if len(cert) == 0 {
		return tls.Certificate{}, crypto.ErrMissingKey
	}

	key, err := cfg.GetKey(fs)
	if err != nil {
		return tls.Certificate{}, err
	}
	if len(key) == 0 {
		return tls.Certificate{}, crypto.ErrMissingKey
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	return pair, nil
}
