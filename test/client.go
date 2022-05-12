package test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gprometheus "github.com/alexfalkowski/go-service/transport/grpc/metrics/prometheus"
	gopentracing "github.com/alexfalkowski/go-service/transport/grpc/trace/opentracing"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	hopentracing "github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewHTTPClient for test.
func NewHTTPClient(lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config) *http.Client {
	return NewHTTPClientWithRoundTripper(lc, logger, cfg, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(lc fx.Lifecycle, logger *zap.Logger, cfg *opentracing.Config, roundTripper http.RoundTripper) *http.Client {
	tracer, _ := hopentracing.NewTracer(hopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	return shttp.NewClient(
		shttp.ClientParams{Version: Version, Config: NewHTTPConfig()},
		shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
		shttp.WithClientMetrics(hprometheus.NewClientMetrics(lc, Version)),
	)
}

// NewGRPCClient for test.
func NewGRPCClient(
	ctx context.Context, lc fx.Lifecycle, gcfg *tgrpc.Config,
	logger *zap.Logger, cfg *opentracing.Config,
	cred credentials.PerRPCCredentials,
) *grpc.ClientConn {
	tracer, _ := gopentracing.NewTracer(gopentracing.TracerParams{Lifecycle: lc, Config: cfg, Version: Version})

	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	if cred != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(cred))
	}

	conn, _ := tgrpc.NewClient(
		tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gcfg.Port), Version: Version, Config: gcfg},
		tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
		tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
		tgrpc.WithClientDialOption(dialOpts...),
		tgrpc.WithClientMetrics(gprometheus.NewClientMetrics(lc, Version)),
	)

	return conn
}
