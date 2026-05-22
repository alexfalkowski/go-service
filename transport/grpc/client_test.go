package grpc_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	transportgrpc.Register(test.FS)

	_, err := transportgrpc.NewClient("none", transportgrpc.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = transportgrpc.NewClient("none", transportgrpc.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
}
