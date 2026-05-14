package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestNewKeyPair(t *testing.T) {
	validTests := []struct {
		config *tls.Config
		name   string
	}{
		{name: "nil config", config: nil},
		{name: "empty config", config: &tls.Config{}},
		{name: "valid config", config: test.NewTLSConfig("certs/cert.pem", "certs/key.pem")},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.HasKeyPair() {
				_, err := tls.NewKeyPair(test.FS, tt.config)
				require.NoError(t, err)
			}
		})
	}

	invalidTests := []struct {
		config        *tls.Config
		name          string
		errMissingKey bool
	}{
		{name: "nil config", errMissingKey: true},
		{name: "empty config", config: &tls.Config{}, errMissingKey: true},
		{name: "missing cert", config: &tls.Config{Key: test.FilePath("certs/key.pem")}, errMissingKey: true},
		{name: "missing key", config: &tls.Config{Cert: test.FilePath("certs/cert.pem")}, errMissingKey: true},
		{name: "empty cert source", config: &tls.Config{Cert: "env:TLS_EMPTY", Key: test.FilePath("certs/key.pem")}, errMissingKey: true},
		{name: "empty key source", config: &tls.Config{Cert: test.FilePath("certs/cert.pem"), Key: "env:TLS_EMPTY"}, errMissingKey: true},
		{name: "invalid key", config: test.NewTLSConfig("certs/client-cert.pem", "secrets/none")},
		{name: "invalid cert", config: test.NewTLSConfig("secrets/none", "certs/client-key.pem")},
		{name: "invalid pair", config: test.NewTLSConfig("secrets/hooks", "certs/client-key.pem")},
	}

	t.Setenv("TLS_EMPTY", "")

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tls.NewKeyPair(test.FS, tt.config)
			if tt.errMissingKey {
				require.ErrorIs(t, err, errors.ErrMissingKey)
				return
			}
			require.Error(t, err)
		})
	}
}
