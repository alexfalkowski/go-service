package transport

import (
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/types/slices"
	"go.uber.org/fx"
)

// ServersParams for transport.
type ServersParams struct {
	fx.In

	HTTP  *http.Server
	GRPC  *grpc.Server
	Debug *debug.Server
}

// NewServers for transport.
func NewServers(params ServersParams) []*server.Service {
	return slices.AppendNotNil([]*server.Service{}, params.HTTP.GetServer(), params.GRPC.GetServer(), params.Debug.GetServer())
}
