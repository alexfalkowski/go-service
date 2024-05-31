package ssh

import (
	"strings"
	"time"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/errors"
	ns "github.com/alexfalkowski/go-service/net/ssh"
	"github.com/alexfalkowski/go-service/server"
	t "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/ssh/handler"
	logger "github.com/alexfalkowski/go-service/transport/ssh/telemetry/logger/zap"
	"github.com/gliderlabs/ssh"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams for SSH.
type ServerParams struct {
	fx.In
	Shutdowner fx.Shutdowner
	Tracer     trace.Tracer
	Meter      metric.Meter
	Config     *Config
	Logger     *zap.Logger
	Handler    handler.Server
	Name       env.Name
	Version    env.Version
}

// Server for SSH.
type Server struct {
	*server.Server
}

// NewServer for SSH.
func NewServer(params ServerParams) (*Server, error) {
	var handler handler.Server = logger.NewServer(params.Logger, params.Handler)

	timeout := timeout(params.Config)
	s := &ssh.Server{
		Banner:  string(params.Name),
		Version: params.Version.String(),
		Handler: func(s ssh.Session) {
			fullCmd := strings.Join(s.Command(), " ")

			err := handler.Handle(s.Context(), s.Command())
			if err != nil {
				s.Write([]byte(fullCmd + ":" + err.Error()))
				s.Exit(1)

				return
			}

			s.Write([]byte(fullCmd + ": successful"))
			s.Exit(0)
		},
		IdleTimeout: timeout,
		MaxTimeout:  timeout,
	}

	c := config(params.Config)

	sv, err := ns.NewServer(s, c)
	if err != nil {
		return nil, errors.Prefix("ssh server", err)
	}

	svr := server.NewServer("ssh", sv, params.Logger, params.Shutdowner)

	return &Server{Server: svr}, nil
}

func config(cfg *Config) *ns.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	c := &ns.Config{
		Port: cfg.GetPort("2222"),
	}

	return c
}

func timeout(cfg *Config) time.Duration {
	if !IsEnabled(cfg) {
		return time.Minute
	}

	return t.MustParseDuration(cfg.Timeout)
}
