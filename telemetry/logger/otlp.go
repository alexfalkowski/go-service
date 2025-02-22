package logger

import (
	"context"
	"log/slog"

	"github.com/alexfalkowski/go-service/errors"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
)

func newOtlpLogger(params Params) (*slog.Logger, error) {
	if err := params.Config.Headers.Secrets(params.FileSystem); err != nil {
		return nil, errors.Prefix("logger", err)
	}

	client, _ := otlp.New(context.Background(), otlp.WithEndpointURL(params.Config.URL), otlp.WithHeaders(params.Config.Headers))
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(params.Name.String()),
		semconv.ServiceVersion(params.Version.String()),
		semconv.DeploymentEnvironmentName(params.Environment.String()),
	)

	provider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(client)), log.WithResource(attrs))
	global.SetLoggerProvider(provider)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = provider.Shutdown(ctx)
			_ = client.Shutdown(ctx)

			return nil
		},
	})

	return otelslog.NewLogger(params.Name.String(), otelslog.WithLoggerProvider(provider)), nil
}
