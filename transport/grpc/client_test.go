package grpc_test

import (
	"testing"

	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net"
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

func TestClientEmptyTLSConfigUsesTLS(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldGRPC())
	_, target := net.ListenNetworkAddress(world.TransportConfig.GRPC.Address)

	conn, err := transportgrpc.NewClient(target, transportgrpc.WithClientTLS(&tls.Config{}))
	require.NoError(t, err)
	defer conn.Close()

	client := v1.NewGreeterServiceClient(conn)
	_, err = client.SayHello(t.Context(), &v1.SayHelloRequest{Name: "test"})
	require.Error(t, err)
}
