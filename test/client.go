package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gtel "github.com/alexfalkowski/go-service/transport/grpc/telemetry"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	htel "github.com/alexfalkowski/go-service/transport/http/telemetry"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewHTTPClient for test.
func NewHTTPClient(lc fx.Lifecycle, logger *zap.Logger, cfg *telemetry.Config, tcfg *transport.Config) *http.Client {
	return NewHTTPClientWithRoundTripper(lc, logger, cfg, tcfg, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(lc fx.Lifecycle, logger *zap.Logger, cfg *telemetry.Config, tcfg *transport.Config, roundTripper http.RoundTripper) *http.Client {
	tracer, _ := htel.NewTracer(htel.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	return shttp.NewClient(
		shttp.ClientParams{Config: &tcfg.HTTP},
		shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
		shttp.WithClientMetrics(htel.NewClientMetrics(lc, Version)),
	)
}

// NewGRPCClient for test.
func NewGRPCClient(
	ctx context.Context, lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *telemetry.Config,
	cred credentials.PerRPCCredentials,
) *grpc.ClientConn {
	tracer, _ := gtel.NewTracer(gtel.TracerParams{Lifecycle: lc, Config: ocfg, Version: Version})

	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	if cred != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(cred))
	}

	conn, _ := tgrpc.NewClient(
		tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", tcfg.Port), Config: &tcfg.GRPC},
		tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
		tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
		tgrpc.WithClientDialOption(dialOpts...),
		tgrpc.WithClientMetrics(gtel.NewClientMetrics(lc, Version)),
	)

	return conn
}
