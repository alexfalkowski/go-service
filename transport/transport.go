package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/soheilhy/cmux"
	"go.uber.org/fx"
)

// RegisterParams for transport.
type RegisterParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Config     *Config
	HTTP       *http.Server
	GRPC       *grpc.Server
}

// Register all the transports.
func Register(params RegisterParams) {
	server := NewServer(params)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			server.Stop(ctx)

			return nil
		},
	})
}

// NewServer for transport.
func NewServer(params RegisterParams) *Server {
	return &Server{sh: params.Shutdowner, cfg: params.Config, http: params.HTTP, grpc: params.GRPC}
}

// Server handles all the transports.
type Server struct {
	mux  cmux.CMux
	sh   fx.Shutdowner
	cfg  *Config
	http *http.Server
	grpc *grpc.Server
}

// Start all the servers.
func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return err
	}

	s.mux = cmux.New(l)
	gl := s.mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	hl := s.mux.Match(cmux.HTTP1())

	go s.grpc.Start(gl)
	go s.http.Start(hl)
	go s.start()

	return nil
}

// Stop all the servers.
func (s *Server) Stop(ctx context.Context) {
	s.grpc.Stop(ctx)
	s.http.Stop(ctx)
	s.mux.Close()
}

func (s *Server) start() error {
	if err := s.mux.Serve(); err != nil {
		return s.sh.Shutdown()
	}

	return nil
}
