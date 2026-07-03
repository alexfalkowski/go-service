package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

func newOtlpLogger(params LoggerParams) (*slog.Logger, error) {
	if err := otlp.ValidateEndpoint(params.Config.GetProtocol(), params.Config.URL, params.Config.Headers); err != nil {
		return nil, err
	}

	exporter, err := newOtlpExporter(params.Config)
	if err != nil {
		return nil, err
	}

	attrs := attributes.NewResource(
		params.Attributes,
		params.ID.String(),
		params.Name.String(),
		params.Version.String(),
		params.Environment.String(),
	)

	provider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(exporter)), log.WithResource(attrs))
	global.SetLoggerProvider(provider)

	params.Lifecycle.Append(di.Hook{
		OnStop: func(ctx context.Context) error {
			// Do not return error as this will stop all others.
			_ = provider.Shutdown(ctx)
			_ = exporter.Shutdown(ctx)

			return nil
		},
	})

	handler := otelslog.NewHandler(params.Name.String(), otelslog.WithLoggerProvider(provider))

	return slog.New(&levelHandler{handler: handler, level: level(params.Config)}), nil
}

func newOtlpExporter(cfg *Config) (log.Exporter, error) {
	switch cfg.GetProtocol() {
	case otlp.ProtocolGRPC:
		opts := []otlploggrpc.Option{otlploggrpc.WithHeaders(cfg.Headers), otlploggrpc.WithInsecure()}
		if !strings.IsEmpty(cfg.URL) {
			opts = append(opts, otlploggrpc.WithEndpoint(cfg.URL))
		}
		return otlploggrpc.New(context.Background(), opts...)
	default:
		opts := []otlploghttp.Option{otlploghttp.WithHeaders(cfg.Headers)}
		if !strings.IsEmpty(cfg.URL) {
			opts = append(opts, otlploghttp.WithEndpointURL(cfg.URL))
		}
		return otlploghttp.New(context.Background(), opts...)
	}
}

type levelHandler struct {
	handler slog.Handler
	level   slog.Level
}

func (h *levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level && h.handler.Enabled(ctx, level)
}

func (h *levelHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelHandler{handler: h.handler.WithAttrs(attrs), level: h.level}
}

func (h *levelHandler) WithGroup(name string) slog.Handler {
	return &levelHandler{handler: h.handler.WithGroup(name), level: h.level}
}
