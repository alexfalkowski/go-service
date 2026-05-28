package transport_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/transport"
	transportgrpc "github.com/alexfalkowski/go-service/v2/transport/grpc"
	transporthttp "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/stretchr/testify/require"
)

func TestNewServers(t *testing.T) {
	require.Empty(t, transport.NewServers(transport.ServersParams{}))

	httpService := newService("http")
	grpcService := newService("grpc")
	debugService := newService("debug")

	services := transport.NewServers(transport.ServersParams{
		HTTP:  &transporthttp.Server{Service: httpService},
		GRPC:  &transportgrpc.Server{Service: grpcService},
		Debug: &debug.Server{Service: debugService},
	})

	require.Equal(t, []*server.Service{httpService, grpcService, debugService}, services)
}

func newService(name string) *server.Service {
	return server.NewService(name, &test.NoopServer{}, nil, test.NewShutdowner())
}
