package metrics

import (
	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// ErrNotFound is returned when the configured metrics reader kind is unknown.
var ErrNotFound = errors.New("metrics: reader not found")

// ReaderParams declares the dependencies required by NewReader.
//
// It is intended for Fx/Dig injection. OTLP/gRPC readers use FS to resolve TLS
// source strings when TLS is configured.
type ReaderParams struct {
	di.In

	// Lifecycle is used to shut down the reader with the application.
	Lifecycle di.Lifecycle

	// Config enables metrics when non-nil and selects the reader/exporter kind.
	Config *Config

	// FS resolves TLS source strings for OTLP/gRPC exporters.
	FS *os.FS

	// Name identifies the service emitting metrics.
	Name env.Name
}

// NewReader constructs an OpenTelemetry SDK metric reader based on Config.
//
// If metrics are disabled, it returns (nil, nil). The constructed reader is
// registered with the provided lifecycle and is shut down on stop. If the
// reader was already shut down, the shutdown error is ignored.
func NewReader(params ReaderParams) (metric.Reader, error) {
	if !params.Config.IsEnabled() {
		return nil, nil
	}

	reader, err := newMetricReader(params.Name, params.FS, params.Config)
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(di.Hook{
		OnStop: func(ctx context.Context) error {
			if err := reader.Shutdown(ctx); err != nil {
				if errors.Is(err, metric.ErrReaderShutdown) {
					return nil
				}
				return err
			}

			return nil
		},
	})
	return reader, nil
}

func newMetricReader(name env.Name, fs *os.FS, cfg *Config) (metric.Reader, error) {
	switch cfg.Kind {
	case "otlp":
		if err := otlp.ValidateEndpoint(otlp.Endpoint{
			Protocol: cfg.GetProtocol(),
			Address:  cfg.URL,
			Headers:  cfg.Headers,
			TLS:      cfg.TLS,
		}); err != nil {
			return nil, prefix(err)
		}

		exporter, err := newOTLPExporter(fs, cfg)
		if err != nil {
			return nil, prefix(err)
		}
		return metric.NewPeriodicReader(exporter, periodicReaderOptions(cfg)...), nil
	case "prometheus":
		exporter, err := prometheus.New(prometheus.WithNamespace(name.String()))
		if err != nil {
			return nil, prefix(err)
		}
		return exporter, nil
	default:
		return nil, ErrNotFound
	}
}

func newOTLPExporter(fs *os.FS, cfg *Config) (metric.Exporter, error) {
	switch cfg.GetProtocol() {
	case otlp.ProtocolGRPC:
		opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithHeaders(cfg.Headers)}
		if cfg.TLS == nil {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		} else {
			conf, err := client.NewConfig(fs, cfg.TLS)
			if err != nil {
				return nil, err
			}
			opts = append(opts, otlpmetricgrpc.WithTLSCredentials(grpc.NewTLS(conf)))
		}
		if !strings.IsEmpty(cfg.URL) {
			opts = append(opts, otlpmetricgrpc.WithEndpoint(cfg.URL))
		}
		return otlpmetricgrpc.New(context.Background(), opts...)
	default:
		opts := []otlpmetrichttp.Option{otlpmetrichttp.WithHeaders(cfg.Headers)}
		if !strings.IsEmpty(cfg.URL) {
			opts = append(opts, otlpmetrichttp.WithEndpointURL(cfg.URL))
		}
		return otlpmetrichttp.New(context.Background(), opts...)
	}
}

func periodicReaderOptions(cfg *Config) []metric.PeriodicReaderOption {
	opts := make([]metric.PeriodicReaderOption, 0, 2)
	if cfg.Interval > 0 {
		opts = append(opts, metric.WithInterval(cfg.Interval.Duration()))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, metric.WithTimeout(cfg.Timeout.Duration()))
	}

	return opts
}
