package tls

import (
	"crypto/tls"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewConfig constructs a `crypto/tls`.Config from cfg.
//
// # Defaults
//
// NewConfig always applies the following defaults:
//   - MinVersion: TLS 1.2 (tls.VersionTLS12)
//   - ClientAuth: require and verify client certificates (mTLS) via tls.RequireAndVerifyClientCert
//
// # Key material loading
//
// If cfg is nil (disabled) or cfg does not have both a certificate and key configured (see cfg.HasKeyPair),
// NewConfig returns a config with the defaults above and does not attempt to load any key material.
//
// When cfg is enabled and has a key pair, certificate and key bytes are resolved via the provided filesystem
// using go-service "source strings" (literal value, `file:` path, or `env:` reference). The resolved PEM
// is then parsed using tls.X509KeyPair and attached to the returned config.
//
// # Errors
//
// NewConfig returns the partially constructed config along with any error encountered while resolving the
// certificate/key sources or parsing them as an X.509 key pair.
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
