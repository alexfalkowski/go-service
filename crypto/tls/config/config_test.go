package config_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestHasKeyMaterial(t *testing.T) {
	tests := []struct {
		config *tls.Config
		name   string
		want   bool
	}{
		{name: "nil"},
		{name: "empty", config: &tls.Config{}},
		{name: "cert", config: &tls.Config{Cert: test.FilePath("certs/cert.pem")}, want: true},
		{name: "key", config: &tls.Config{Key: test.FilePath("certs/key.pem")}, want: true},
		{name: "ca", config: &tls.Config{CA: test.FilePath("certs/rootCA.pem")}},
		{name: "server name", config: &tls.Config{ServerName: "localhost"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.config.HasKeyMaterial())
		})
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		config *tls.Config
		name   string
		want   bool
	}{
		{name: "nil"},
		{name: "empty", config: &tls.Config{}},
		{name: "cert", config: &tls.Config{Cert: test.FilePath("certs/cert.pem")}, want: true},
		{name: "key", config: &tls.Config{Key: test.FilePath("certs/key.pem")}, want: true},
		{name: "ca", config: &tls.Config{CA: test.FilePath("certs/rootCA.pem")}, want: true},
		{name: "server name", config: &tls.Config{ServerName: "localhost"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.config.IsEnabled())
		})
	}
}
