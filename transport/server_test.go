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
	httpService := newService("http")
	grpcService := newService("grpc")
	debugService := newService("debug")

	tests := []struct {
		name   string
		params transport.ServersParams
		want   []*server.Service
	}{
		{name: "none", want: []*server.Service{}},
		{
			name: "http only",
			params: transport.ServersParams{
				HTTP: &transporthttp.Server{Service: httpService},
			},
			want: []*server.Service{httpService},
		},
		{
			name: "grpc only",
			params: transport.ServersParams{
				GRPC: &transportgrpc.Server{Service: grpcService},
			},
			want: []*server.Service{grpcService},
		},
		{
			name: "debug only",
			params: transport.ServersParams{
				Debug: &debug.Server{Service: debugService},
			},
			want: []*server.Service{debugService},
		},
		{
			name: "http and debug",
			params: transport.ServersParams{
				HTTP:  &transporthttp.Server{Service: httpService},
				Debug: &debug.Server{Service: debugService},
			},
			want: []*server.Service{httpService, debugService},
		},
		{
			name: "all",
			params: transport.ServersParams{
				HTTP:  &transporthttp.Server{Service: httpService},
				GRPC:  &transportgrpc.Server{Service: grpcService},
				Debug: &debug.Server{Service: debugService},
			},
			want: []*server.Service{httpService, grpcService, debugService},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, transport.NewServers(tt.params))
		})
	}
}

func newService(name string) *server.Service {
	return server.NewService(name, &test.NoopServer{}, nil, test.NewShutdowner())
}
