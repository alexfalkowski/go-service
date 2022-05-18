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
	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

// Register all the transports.
func Register(params RegisterParams) {
	var mux cmux.CMux

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			l, err := net.Listen("tcp", fmt.Sprintf(":%s", params.Config.Port))
			if err != nil {
				return err
			}

			mux = cmux.New(l)
			gl := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
			hl := mux.Match(cmux.HTTP1())

			go params.GRPCServer.Start(gl)
			go params.HTTPServer.Start(hl)
			go start(mux, params.Shutdowner)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.GRPCServer.Stop(ctx)
			params.HTTPServer.Stop(ctx)
			mux.Close()

			return nil
		},
	})
}

func start(mux cmux.CMux, sh fx.Shutdowner) error {
	if err := mux.Serve(); err != nil {
		return sh.Shutdown()
	}

	return nil
}
