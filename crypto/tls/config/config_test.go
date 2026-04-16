package config_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	validTests := []struct {
		config *tls.Config
		name   string
	}{
		{name: "nil config", config: nil},
		{name: "empty config", config: &tls.Config{}},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tls.NewConfig(test.FS, tt.config)
			require.NoError(t, err)
			require.NotNil(t, config)
		})
	}

	invalidTests := []struct {
		config *tls.Config
		name   string
	}{
		{name: "invalid key", config: test.NewTLSConfig("certs/client-cert.pem", "secrets/none")},
		{name: "invalid cert", config: test.NewTLSConfig("secrets/none", "certs/client-key.pem")},
		{name: "invalid pair", config: test.NewTLSConfig("secrets/hooks", "certs/client-key.pem")},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tls.NewConfig(test.FS, tt.config)
			require.Error(t, err)
		})
	}
}
