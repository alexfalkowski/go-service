package client_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	require.False(t, (*client.Config)(nil).IsEnabled())
	require.True(t, (&client.Config{}).IsEnabled())
}

func TestNewConfig(t *testing.T) {
	tlsConfig, err := client.NewConfig(test.FS, test.NewTLSClientConfig())
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
	require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	require.Len(t, tlsConfig.Certificates, 1)
	require.NotNil(t, tlsConfig.RootCAs)
	require.Equal(t, "localhost", tlsConfig.ServerName)

	data, err := test.FS.ReadSource(test.FilePath("certs/cert.pem"))
	require.NoError(t, err)

	block, _ := pem.Decode(data)
	require.NotNil(t, block)

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	_, err = cert.Verify(x509.VerifyOptions{DNSName: tlsConfig.ServerName, Roots: tlsConfig.RootCAs})
	require.NoError(t, err)
}

func TestNewConfigDefaults(t *testing.T) {
	tests := []struct {
		config *config.Config
		name   string
	}{
		{name: "nil"},
		{name: "empty", config: &config.Config{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, err := client.NewConfig(test.FS, tt.config)
			require.NoError(t, err)
			require.NotNil(t, tlsConfig)
			require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
			require.Empty(t, tlsConfig.Certificates)
			require.Nil(t, tlsConfig.RootCAs)
			require.Empty(t, tlsConfig.ServerName)
		})
	}
}

func TestNewConfigWithKeyPairOnly(t *testing.T) {
	tlsConfig, err := client.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/client-cert.pem"),
		Key:  test.FilePath("certs/client-key.pem"),
	})
	require.NoError(t, err)
	require.Len(t, tlsConfig.Certificates, 1)
	require.Nil(t, tlsConfig.RootCAs)
}

func TestNewConfigWithCAOnly(t *testing.T) {
	tlsConfig, err := client.NewConfig(test.FS, &config.Config{CA: test.FilePath("certs/rootCA.pem")})
	require.NoError(t, err)
	require.Empty(t, tlsConfig.Certificates)
	require.NotNil(t, tlsConfig.RootCAs)
}

func TestNewConfigInvalidKeyPair(t *testing.T) {
	_, err := client.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/client-cert.pem"),
		Key:  test.FilePath("secrets/none"),
	})
	require.Error(t, err)
}

func TestNewConfigInvalidCA(t *testing.T) {
	_, err := client.NewConfig(test.FS, &config.Config{CA: "invalid ca"})
	require.Error(t, err)
	require.True(t, errors.Is(err, config.ErrInvalidCA))
}
