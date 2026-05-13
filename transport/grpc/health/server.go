package health

import (
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-health/v2/subscriber"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ServerParams defines dependencies for constructing the gRPC health `Server` implementation.
//
// It is an Fx parameter struct (`di.In`) that provides the underlying `*health.Server`
// from `github.com/alexfalkowski/go-health/v2/server`, which maintains health observers.
type ServerParams struct {
	di.In

	// Server is the underlying health server that stores and exposes health observers.
	//
	// It is expected to be non-nil when this gRPC health service is wired.
	Server *server.Server
}

// NewServer constructs a new gRPC health `Server` implementation.
//
// The returned server implements the standard gRPC health service
// (`grpc.health.v1.Health`) and delegates health state lookups to the provided
// underlying `*health.Server`.
func NewServer(params ServerParams) *Server {
	return &Server{server: params.Server}
}

// Server implements the standard gRPC health service.
//
// It exposes health state by querying observers registered in the underlying `*health.Server`.
// The service is typically used by load balancers and orchestration systems to determine whether
// a server is serving traffic.
type Server struct {
	health.UnimplementedServer
	server *server.Server
}

// Check returns the health status for a single service.
//
// The requested service name is taken from `req.GetService()`. The health state is resolved by looking up
// an observer from the underlying `*health.Server` using the "grpc" transport kind.
//
// Error mapping:
//   - If the requested service does not exist, it returns `codes.NotFound`.
//   - Otherwise, it returns a `Response` whose status is derived from the observer's error state.
func (s *Server) Check(_ context.Context, req *health.Request) (*health.Response, error) {
	observer, err := s.server.Observer(req.GetService(), "grpc")
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &health.Response{Status: s.status(observer)}, nil
}

// List returns the health status for all registered services.
//
// It enumerates all observers registered for the "grpc" transport kind and returns their current serving
// statuses.
func (s *Server) List(_ context.Context, req *health.ListRequest) (*health.ListResponse, error) {
	res := &health.ListResponse{Statuses: map[string]*health.Response{}}
	for name, observer := range s.server.Observers("grpc") {
		res.Statuses[name] = &health.Response{Status: s.status(observer)}
	}

	return res, nil
}

// Watch streams health status updates for a single service until the client cancels.
//
// The initial status is sent immediately. When the requested service is unknown, Watch sends
// `SERVICE_UNKNOWN` and keeps the stream open so clients can observe the service becoming available later.
//
// This package's underlying health server exposes observer state but not a push-based watch API, so Watch
// polls that in-memory state and only emits a response when the effective serving status changes.
func (s *Server) Watch(req *health.Request, w health.WatchServer) error {
	service := req.GetService()
	current := health.Unknown
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		next := s.watchStatus(service)
		if next != current {
			current = next
			if err := w.Send(&health.Response{Status: current}); err != nil {
				return status.Error(codes.Canceled, "stream has ended")
			}
		}

		select {
		case <-w.Context().Done():
			return status.Error(codes.Canceled, w.Context().Err().Error())
		case <-ticker.C:
		}
	}
}

func (s *Server) status(observer *subscriber.Observer) health.Status {
	if err := observer.Error(); err != nil {
		return health.NotServing
	}

	return health.Serving
}

func (s *Server) watchStatus(service string) health.Status {
	observer, err := s.server.Observer(service, "grpc")
	if err != nil {
		return health.ServiceUnknown
	}

	return s.status(observer)
}
