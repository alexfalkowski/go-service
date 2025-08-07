package tls

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewConfig for tls.
func NewConfig(fs *os.FS, cfg *Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !cfg.IsEnabled() || !cfg.HasKeyPair() {
		return config, nil
	}

	cert, err := cfg.GetCert(fs)
	if err != nil {
		return config, err
	}

	key, err := cfg.GetKey(fs)
	if err != nil {
		return config, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return config, err
	}

	config.Certificates = []tls.Certificate{pair}
	return config, nil
}
