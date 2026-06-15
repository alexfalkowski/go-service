package health

import (
	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-health/v2/watcher"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ServerParams defines dependencies for constructing the gRPC health [Server] implementation.
//
// It is an Fx parameter struct ([di.In]) that provides the underlying *[health.Server]
// from [github.com/alexfalkowski/go-health/v2/server], which maintains health observers.
type ServerParams struct {
	di.In

	// Server is the underlying health server that stores and exposes health observers.
	//
	// It is expected to be non-nil when this gRPC health service is wired.
	Server *server.Server
}

// NewServer constructs a new gRPC health [Server] implementation.
//
// The returned server implements the standard gRPC health service
// (`grpc.health.v1.Health`) and delegates health state lookups to the provided
// underlying *[health.Server].
func NewServer(params ServerParams) *Server {
	return &Server{server: params.Server}
}

// Server implements the standard gRPC health service.
//
// It exposes health state by querying observers registered in the underlying *[health.Server].
// The service is typically used by load balancers and orchestration systems to determine whether
// a server is serving traffic.
type Server struct {
	health.UnimplementedServer
	server *server.Server
}

// Check returns the health status for a single service.
//
// The requested service name is taken from `req.GetService()`. The health state is resolved by looking up
// an observer from the underlying *[health.Server] using the "grpc" transport kind.
//
// Error mapping:
//   - If the requested service does not exist, it returns [codes.NotFound].
//   - Otherwise, it returns a `Response` whose status is derived from the observer's error state.
func (s *Server) Check(_ context.Context, req *health.Request) (*health.Response, error) {
	service := req.GetService()
	if strings.IsEmpty(service) {
		return &health.Response{Status: healthStatus(s.server.Error("grpc"))}, nil
	}

	observer, err := s.server.Observer(service, "grpc")
	if err != nil {
		return nil, status.SafeError(codes.NotFound, err)
	}

	return &health.Response{Status: healthStatus(observer.Error())}, nil
}

// List returns the health status for all registered services.
//
// It enumerates all observers registered for the "grpc" transport kind and returns their current serving
// statuses.
func (s *Server) List(_ context.Context, _ *health.ListRequest) (*health.ListResponse, error) {
	res := &health.ListResponse{Statuses: map[string]*health.Response{}}
	for name, observer := range s.server.Observers("grpc") {
		res.Statuses[name] = &health.Response{Status: healthStatus(observer.Error())}
	}

	return res, nil
}

// Watch streams health status updates for a single service until the client cancels.
//
// The initial status is sent immediately. When the requested service is unknown,
// Watch sends `SERVICE_UNKNOWN` and keeps the stream open until the client cancels.
func (s *Server) Watch(req *health.Request, w health.WatchServer) error {
	service := req.GetService()
	if strings.IsEmpty(service) {
		return s.watch(w, s.server.Watch("grpc"))
	}

	observer, err := s.server.Observer(service, "grpc")
	if err != nil {
		return sendUnknownStatus(w)
	}

	return s.watch(w, observer.Watch())
}

func (s *Server) watch(w health.WatchServer, sub watcher.Subscription) error {
	defer sub.Close()

	current := health.Unknown
	for {
		select {
		case err, ok := <-sub.Receive():
			if !ok {
				return status.Error(codes.Canceled, "watch has ended")
			}

			next := healthStatus(err)
			if next != current {
				current = next
				if err := sendStatus(w, current); err != nil {
					return err
				}
			}
		case <-w.Context().Done():
			return status.SafeError(codes.Canceled, w.Context().Err())
		}
	}
}

func healthStatus(err error) health.Status {
	if err != nil {
		return health.NotServing
	}

	return health.Serving
}

func sendUnknownStatus(w health.WatchServer) error {
	if err := sendStatus(w, health.ServiceUnknown); err != nil {
		return err
	}

	<-w.Context().Done()

	return status.SafeError(codes.Canceled, w.Context().Err())
}

func sendStatus(w health.WatchServer, st health.Status) error {
	if err := w.Send(&health.Response{Status: st}); err != nil {
		if _, ok := status.FromError(err); ok {
			return err
		}

		return status.Error(codes.Canceled, "stream has ended")
	}

	return nil
}
