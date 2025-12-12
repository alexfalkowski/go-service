package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	grpc.Register(test.FS)

	_, err := grpc.NewClient("none", grpc.WithClientTLS(&tls.Config{Cert: "bob", Key: "bob"}))
	require.Error(t, err)

	_, err = grpc.NewClient("none", grpc.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
}
