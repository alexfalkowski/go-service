package server_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetMaxReceiveSizeReturnsConfiguredOrDefaultValue(t *testing.T) {
	tests := []struct {
		cfg  *server.Config
		name string
		want bytes.Size
	}{
		{name: "nil", want: bytes.DefaultSize},
		{name: "zero", cfg: &server.Config{}, want: bytes.DefaultSize},
		{name: "explicit", cfg: &server.Config{MaxReceiveSize: 64}, want: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.cfg.GetMaxReceiveSize())
		})
	}
}

func TestConfigIsEnabledUnlessNil(t *testing.T) {
	require.False(t, (*server.Config)(nil).IsEnabled())
	require.True(t, (&server.Config{}).IsEnabled())
}

func TestNewConfigLoadsMutuallyAuthenticatedTLS(t *testing.T) {
	tlsConfig, err := server.NewConfig(test.FS, test.NewTLSServerConfig())
	require.NoError(t, err)
	require.NotNil(t, tlsConfig)
	require.Equal(t, uint16(tls.VersionTLS13), tlsConfig.MinVersion)
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

func TestNewConfigAcceptsOptionalTLSMaterial(t *testing.T) {
	tests := []struct {
		config       *config.Config
		name         string
		certificates int
	}{
		{name: "nil"},
		{name: "empty", config: &config.Config{}},
		{name: "server name", config: &config.Config{ServerName: "localhost"}},
		{
			name: "key pair only",
			config: &config.Config{
				Cert: test.FilePath("certs/cert.pem"),
				Key:  test.FilePath("certs/key.pem"),
			},
			certificates: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, err := server.NewConfig(test.FS, tt.config)
			require.NoError(t, err)
			require.NotNil(t, tlsConfig)
			require.Equal(t, uint16(tls.VersionTLS13), tlsConfig.MinVersion)
			require.Len(t, tlsConfig.Certificates, tt.certificates)
			require.Nil(t, tlsConfig.ClientCAs)
			require.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
		})
	}
}

func TestNewConfigRejectsIncompleteKeyPair(t *testing.T) {
	tests := []struct {
		config *config.Config
		name   string
	}{
		{name: "CA only", config: &config.Config{CA: test.FilePath("certs/rootCA.pem")}},
		{name: "cert only", config: &config.Config{Cert: test.FilePath("certs/cert.pem")}},
		{name: "key only", config: &config.Config{Key: test.FilePath("certs/key.pem")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsConfig, err := server.NewConfig(test.FS, tt.config)
			require.ErrorIs(t, err, server.ErrMissingKeyPair)
			require.Empty(t, tlsConfig.Certificates)
			require.Nil(t, tlsConfig.ClientCAs)
			require.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)
		})
	}
}

func TestNewConfigRejectsInvalidKeyPair(t *testing.T) {
	_, err := server.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/cert.pem"),
		Key:  test.FilePath("secrets/none"),
	})
	require.Error(t, err)
}

func TestNewConfigRejectsInvalidCA(t *testing.T) {
	_, err := server.NewConfig(test.FS, &config.Config{
		Cert: test.FilePath("certs/cert.pem"),
		Key:  test.FilePath("certs/key.pem"),
		CA:   "invalid ca",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, config.ErrInvalidCA)
}
