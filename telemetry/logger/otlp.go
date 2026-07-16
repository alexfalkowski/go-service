package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
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
	if err := otlp.ValidateEndpoint(otlp.Endpoint{
		Protocol: params.Config.GetProtocol(),
		Address:  params.Config.URL,
		Headers:  params.Config.Headers,
		TLS:      params.Config.TLS,
	}); err != nil {
		return nil, err
	}

	exporter, err := newOtlpExporter(params)
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

	provider := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(exporter, batchProcessorOptions(params.Config)...)), log.WithResource(attrs))
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

func newOtlpExporter(params LoggerParams) (log.Exporter, error) {
	switch params.Config.GetProtocol() {
	case otlp.ProtocolGRPC:
		opts := []otlploggrpc.Option{otlploggrpc.WithHeaders(params.Config.Headers)}
		if params.Config.TLS == nil {
			opts = append(opts, otlploggrpc.WithInsecure())
		} else {
			conf, err := client.NewConfig(params.FS, params.Config.TLS)
			if err != nil {
				return nil, err
			}
			opts = append(opts, otlploggrpc.WithTLSCredentials(grpc.NewTLS(conf)))
		}
		if !strings.IsEmpty(params.Config.URL) {
			opts = append(opts, otlploggrpc.WithEndpoint(params.Config.URL))
		}
		return otlploggrpc.New(context.Background(), opts...)
	default:
		opts := []otlploghttp.Option{otlploghttp.WithHeaders(params.Config.Headers)}
		if !strings.IsEmpty(params.Config.URL) {
			opts = append(opts, otlploghttp.WithEndpointURL(params.Config.URL))
		}
		return otlploghttp.New(context.Background(), opts...)
	}
}

func batchProcessorOptions(cfg *Config) []log.BatchProcessorOption {
	opts := make([]log.BatchProcessorOption, 0, 4)
	if cfg.BatchTimeout > 0 {
		opts = append(opts, log.WithExportInterval(cfg.BatchTimeout.Duration()))
	}
	if cfg.ExportTimeout > 0 {
		opts = append(opts, log.WithExportTimeout(cfg.ExportTimeout.Duration()))
	}
	if cfg.MaxQueueSize > 0 {
		opts = append(opts, log.WithMaxQueueSize(cfg.MaxQueueSize))
	}
	if cfg.MaxExportBatchSize > 0 {
		opts = append(opts, log.WithExportMaxBatchSize(cfg.MaxExportBatchSize))
	}

	return opts
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
