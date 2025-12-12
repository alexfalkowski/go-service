package tls_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configs := []*tls.Config{nil, {}}

	for _, c := range configs {
		c, err := tls.NewConfig(test.FS, c)
		require.NoError(t, err)
		require.NotNil(t, c)
	}

	configs = []*tls.Config{
		test.NewTLSConfig("certs/client-cert.pem", "secrets/none"),
		test.NewTLSConfig("secrets/none", "certs/client-key.pem"),
		test.NewTLSConfig("secrets/hooks", "certs/client-key.pem"),
	}

	for _, c := range configs {
		_, err := tls.NewConfig(test.FS, c)
		require.Error(t, err)
	}
}
