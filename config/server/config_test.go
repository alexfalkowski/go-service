package server_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetMaxReceiveSize(t *testing.T) {
	tests := []struct {
		cfg  *server.Config
		name string
		want bytes.Size
	}{
		{name: "nil", want: server.DefaultMaxReceiveSize},
		{name: "zero", cfg: &server.Config{}, want: server.DefaultMaxReceiveSize},
		{name: "explicit", cfg: &server.Config{MaxReceiveSize: 64}, want: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetMaxReceiveSize())
		})
	}
}

func TestConfigRejectsNegativeMaxReceiveSize(t *testing.T) {
	cfg := &server.Config{MaxReceiveSize: -1}
	require.Error(t, test.Validator.Struct(cfg))
}

func TestConfig(t *testing.T) {
	require.False(t, (*server.Config)(nil).IsEnabled())
	require.True(t, (&server.Config{}).IsEnabled())
}

func TestNewConfig(t *testing.T) {
	tlsConfig, err := server.NewConfig(test.FS, test.NewTLSServerConfig())
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
	require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	require.Len(t, tlsConfig.Certificates, 1)
	require.NotNil(t, tlsConfig.ClientCAs)
	require.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)

	data, err := test.FS.ReadSource(test.FilePath("certs/client-cert.pem"))
	require.NoError(t, err)

	block, _ := pem.Decode(data)
	require.NotNil(t, block)

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	_, err = cert.Verify(x509.VerifyOptions{Roots: tlsConfig.ClientCAs})
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
			tlsConfig, err := server.NewConfig(test.FS, tt.config)
			require.NoError(t, err)
			require.NotNil(t, tlsConfig)
			require.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
			require.Empty(t, tlsConfig.Certificates)
			require.Nil(t, tlsConfig.ClientCAs)
			require.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
		})
	}
}

func TestNewConfigWithKeyPairOnly(t *testing.T) {
	tlsConfig, err := server.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/cert.pem"),
		Key:  test.FilePath("certs/key.pem"),
	})
	require.NoError(t, err)
	require.Len(t, tlsConfig.Certificates, 1)
	require.Nil(t, tlsConfig.ClientCAs)
	require.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
}

func TestNewConfigWithCAOnly(t *testing.T) {
	tlsConfig, err := server.NewConfig(test.FS, &config.Config{CA: test.FilePath("certs/rootCA.pem")})
	require.NoError(t, err)
	require.Empty(t, tlsConfig.Certificates)
	require.NotNil(t, tlsConfig.ClientCAs)
	require.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
}

func TestNewConfigInvalidKeyPair(t *testing.T) {
	_, err := server.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/cert.pem"),
		Key:  test.FilePath("secrets/none"),
	})
	require.Error(t, err)
}

func TestNewConfigInvalidCA(t *testing.T) {
	_, err := server.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/cert.pem"),
		Key:  test.FilePath("certs/key.pem"),
		CA:   "invalid ca",
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, config.ErrInvalidCA))
}
