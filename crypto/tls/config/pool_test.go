package config_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestNewCertPool(t *testing.T) {
	pool, err := tls.NewCertPool(test.FS, test.NewTLSServerConfig())
	require.NoError(t, err)
	require.NotNil(t, pool)

	data, err := test.FS.ReadSource(test.FilePath("certs/client-cert.pem"))
	require.NoError(t, err)

	block, _ := pem.Decode(data)
	require.NotNil(t, block)

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	_, err = cert.Verify(x509.VerifyOptions{Roots: pool})
	require.NoError(t, err)
}

func TestNewCertPoolInvalid(t *testing.T) {
	_, err := tls.NewCertPool(test.FS, &tls.Config{CA: "invalid ca"})
	require.Error(t, err)
	require.True(t, errors.Is(err, tls.ErrInvalidCA))
}

func TestNewCertPoolSourceError(t *testing.T) {
	_, err := tls.NewCertPool(test.FS, &tls.Config{CA: test.FilePath("certs/missing.pem")})
	require.Error(t, err)
	require.False(t, errors.Is(err, tls.ErrInvalidCA))
}
