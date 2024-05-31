package zap

import (
	tz "github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/transport/ssh/handler"
	"github.com/gliderlabs/ssh"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewServer for SSH.
func NewServer(logger *zap.Logger, handler handler.Server) *Server {
	return &Server{logger: logger, handler: handler}
}

// Server for SSH.
type Server struct {
	logger  *zap.Logger
	handler handler.Server
}

// Handle session.
func (s *Server) Handle(ctx ssh.Context, cmd []string) error {
	fields := []zapcore.Field{
		zap.Strings("command", cmd),
		zap.String("clientVersion", ctx.ClientVersion()),
	}

	err := s.handler.Handle(ctx, cmd)
	tz.LogWithLogger(message("received session"), err, s.logger, fields...)

	return err
}

// NewClient for ssh.
func NewClient(logger *zap.Logger, runner handler.Client) *Client {
	return &Client{logger: logger, runner: runner}
}

// Client for SSH.
type Client struct {
	logger *zap.Logger
	runner handler.Client
}

// Run command.
func (r *Client) Run(cmd string) ([]byte, error) {
	fields := []zapcore.Field{
		zap.String("command", cmd),
	}

	b, err := r.runner.Run(cmd)
	tz.LogWithLogger(message("run command"), err, r.logger, fields...)

	return b, err
}

func message(msg string) string {
	return "ssh: " + msg
}
