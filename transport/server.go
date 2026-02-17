package transport

import (
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/types/slices"
)

// ServersParams defines dependencies used to collect transport servers.
type ServersParams struct {
	di.In
	HTTP  *http.Server
	GRPC  *grpc.Server
	Debug *debug.Server
}

// NewServers returns the non-nil transport services.
func NewServers(params ServersParams) []*server.Service {
	return slices.AppendNotNil([]*server.Service{}, params.HTTP.GetService(), params.GRPC.GetService(), params.Debug.GetService())
}
