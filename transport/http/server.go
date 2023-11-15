package http

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/cors"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	szap "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
)

// ServerParams for HTTP.
type ServerParams struct {
	fx.In

	Shutdowner fx.Shutdowner
	Config     *Config
	Logger     *zap.Logger
	Tracer     tracer.Tracer
	Meter      metric.Meter
}

// Server for HTTP.
type Server struct {
	Mux    *runtime.ServeMux
	server *http.Server
	params ServerParams
}

// NewServer for HTTP.
func NewServer(params ServerParams) (*Server, error) {
	opts := []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(customMatcher),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}
	mux := runtime.NewServeMux(opts...)

	var handler http.Handler = mux

	handler = cors.New().Handler(handler)

	h, err := metrics.NewHandler(params.Meter, handler)
	if err != nil {
		return nil, err
	}

	handler = h
	handler = tracer.NewHandler(params.Tracer, handler)
	handler = szap.NewHandler(params.Logger, handler)
	handler = meta.NewHandler(handler)

	s := &http.Server{
		Handler:           handler,
		ReadTimeout:       time.Timeout,
		WriteTimeout:      time.Timeout,
		IdleTimeout:       time.Timeout,
		ReadHeaderTimeout: time.Timeout,
	}

	server := &Server{
		Mux:    mux,
		server: s,
		params: params,
	}

	return server, nil
}

// Start the server.
func (s *Server) Start(listener net.Listener) {
	s.params.Logger.Info("starting http server", zap.String("addr", listener.Addr().String()))

	if err := s.serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fields := []zapcore.Field{zap.String("addr", listener.Addr().String()), zap.Error(err)}

		if err := s.params.Shutdowner.Shutdown(); err != nil {
			fields = append(fields, zap.NamedError("shutdown_error", err))
		}

		s.params.Logger.Error("could not start http server", fields...)
	}
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) {
	message := "stopping http server"
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.params.Logger.Error(message, zap.Error(err))
	} else {
		s.params.Logger.Info(message)
	}
}

func (s *Server) serve(l net.Listener) error {
	if s.params.Config.Security.IsEnabled() {
		return s.server.ServeTLS(l, s.params.Config.Security.CertFile, s.params.Config.Security.KeyFile)
	}

	return s.server.Serve(l)
}

func customMatcher(key string) (string, bool) {
	switch key {
	case "Request-Id", "Geolocation", "X-Forwarded-For":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
