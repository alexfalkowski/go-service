package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	http.Register(test.FS)

	_, err := http.NewClient(http.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = http.NewClient(http.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
}
