package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/transport"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gotel "github.com/alexfalkowski/go-service/transport/grpc/otel"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hotel "github.com/alexfalkowski/go-service/transport/http/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewHTTPClient for test.
func NewHTTPClient(lc fx.Lifecycle, logger *zap.Logger, cfg *otel.Config, tcfg *transport.Config) *http.Client {
	return NewHTTPClientWithRoundTripper(lc, logger, cfg, tcfg, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(lc fx.Lifecycle, logger *zap.Logger, cfg *otel.Config, tcfg *transport.Config, roundTripper http.RoundTripper) *http.Client {
	tracer, _ := hotel.NewTracer(hotel.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	return shttp.NewClient(
		shttp.ClientParams{Config: &tcfg.HTTP},
		shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
		shttp.WithClientMetrics(hprometheus.NewClientMetrics(lc, Version)),
	)
}

// NewGRPCClient for test.
func NewGRPCClient(
	ctx context.Context, lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *otel.Config,
	cred credentials.PerRPCCredentials,
) *grpc.ClientConn {
	tracer, _ := gotel.NewTracer(gotel.TracerParams{Lifecycle: lc, Config: ocfg, Version: Version})

	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	if cred != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(cred))
	}

	conn, _ := tgrpc.NewClient(
		tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", tcfg.Port), Config: &tcfg.GRPC},
		tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
		tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
		tgrpc.WithClientDialOption(dialOpts...),
		tgrpc.WithClientMetrics(gprometheus.NewClientMetrics(lc, Version)),
	)

	return conn
}
