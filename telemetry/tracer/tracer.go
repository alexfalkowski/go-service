package tracer

import (
	"github.com/alexfalkowski/go-service/v2/config/client"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-sync"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	enabled      sync.Bool
	noopProvider trace.TracerProvider = noop.NewTracerProvider()
)

// ReadOnlySpan is an alias for [go.opentelemetry.io/otel/sdk/trace.ReadOnlySpan].
type ReadOnlySpan = sdk.ReadOnlySpan

// SpanExporter is an alias for [go.opentelemetry.io/otel/sdk/trace.SpanExporter].
type SpanExporter = sdk.SpanExporter

// Provider is an alias for [trace.TracerProvider].
type Provider = trace.TracerProvider

// SpanContext is an alias for [trace.SpanContext].
type SpanContext = trace.SpanContext

// SpanContextConfig is an alias for [trace.SpanContextConfig].
type SpanContextConfig = trace.SpanContextConfig

// SpanID is an alias for [trace.SpanID].
type SpanID = trace.SpanID

// SDKProvider is an alias for [go.opentelemetry.io/otel/sdk/trace.TracerProvider].
type SDKProvider = sdk.TracerProvider

// TraceFlags is an alias for [trace.TraceFlags].
type TraceFlags = trace.TraceFlags

// TraceID is an alias for [trace.TraceID].
type TraceID = trace.TraceID

// ProviderOption is an alias for [go.opentelemetry.io/otel/sdk/trace.TracerProviderOption].
type ProviderOption = sdk.TracerProviderOption

// FlagsSampled is an alias for [trace.FlagsSampled].
const FlagsSampled = trace.FlagsSampled

// ContextWithRemoteSpanContext returns a copy of parent with sc set as the
// current remote span context.
func ContextWithRemoteSpanContext(parent context.Context, sc SpanContext) context.Context {
	return trace.ContextWithRemoteSpanContext(parent, sc)
}

// NewSpanContext constructs a SpanContext from config.
func NewSpanContext(config SpanContextConfig) SpanContext {
	return trace.NewSpanContext(config)
}

// NewProvider constructs an OpenTelemetry SDK tracer provider.
func NewProvider(opts ...ProviderOption) *SDKProvider {
	return sdk.NewTracerProvider(opts...)
}

// WithSyncer configures a tracer provider to export spans synchronously.
func WithSyncer(exporter SpanExporter) ProviderOption {
	return sdk.WithSyncer(exporter)
}

// NewNoopProvider constructs a no-op tracer provider.
func NewNoopProvider() Provider {
	return noop.NewTracerProvider()
}

// TracerParams declares the dependencies required by Register.
//
// It is intended for Fx/Dig injection and includes service identity fields used to
// populate OpenTelemetry resource attributes.
type TracerParams struct {
	di.In

	// Lifecycle is used to start and stop the exporter/provider with the application.
	Lifecycle di.Lifecycle

	// Config enables tracing when non-nil and supplies exporter settings.
	Config *Config

	// FS resolves TLS source strings for OTLP/gRPC exporters.
	FS *os.FS

	// Attributes are optional OpenTelemetry resource attributes attached to traces.
	Attributes attributes.Map `optional:"true"`

	// ID is the host identifier used for the resource's host.id attribute.
	ID env.ID

	// Name is the service name used for the resource's service.name attribute.
	Name env.Name

	// Version is the service version used for the resource's service.version attribute.
	Version env.Version

	// Environment is the deployment environment name used for the resource's
	// deployment.environment.name attribute.
	Environment env.Environment
}

// Register configures and installs a global OpenTelemetry tracer provider.
//
// When tracing is configured with kind "otlp", Register:
//
//  1. Creates an OTLP trace exporter using [Config.Protocol], [Config.URL], and [Config.Headers].
//  2. Creates an OpenTelemetry resource describing the running service.
//  3. Installs a [go.opentelemetry.io/otel/sdk/trace.TracerProvider] globally.
//  4. Appends lifecycle hooks to start the exporter on application start and to
//     shut down the provider/exporter on application stop.
//
// If Config is nil or Kind is empty, Register installs the noop provider. Unknown
// non-empty kinds return ErrNotFound.
//
// Shutdown errors are intentionally ignored to avoid blocking other stop hooks.
func Register(params TracerParams) error {
	if !params.Config.IsEnabled() {
		setProvider(noopProvider, false)
		return nil
	}

	switch params.Config.Kind {
	case "otlp":
		if err := otlp.ValidateEndpoint(otlp.Endpoint{
			Protocol: params.Config.GetProtocol(),
			Address:  params.Config.URL,
			Headers:  params.Config.Headers,
			TLS:      params.Config.TLS,
		}); err != nil {
			return prefix(err)
		}

		exporter, err := newOTLPExporter(params)
		if err != nil {
			return prefix(err)
		}
		attrs := attributes.NewResource(
			params.Attributes,
			params.ID.String(),
			params.Name.String(),
			params.Version.String(),
			params.Environment.String(),
		)

		providerOpts := []sdk.TracerProviderOption{sdk.WithResource(attrs), sdk.WithBatcher(exporter)}
		sampler, err := newSampler(params.Config.Sampler)
		if err != nil {
			return prefix(err)
		}
		if sampler != nil {
			providerOpts = append(providerOpts, sdk.WithSampler(sampler))
		}

		provider := sdk.NewTracerProvider(providerOpts...)
		setProvider(provider, true)

		params.Lifecycle.Append(di.Hook{
			OnStart: func(ctx context.Context) error {
				return prefix(exporter.Start(ctx))
			},
			OnStop: func(ctx context.Context) error {
				// Do not return error as this will stop all others.
				_ = provider.Shutdown(ctx)
				_ = exporter.Shutdown(ctx)
				setProvider(noopProvider, false)

				return nil
			},
		})

		return nil
	default:
		return ErrNotFound
	}
}

func newOTLPExporter(params TracerParams) (*otlptrace.Exporter, error) {
	switch params.Config.GetProtocol() {
	case otlp.ProtocolGRPC:
		opts := []otlptracegrpc.Option{otlptracegrpc.WithHeaders(params.Config.Headers)}
		if params.Config.TLS == nil {
			opts = append(opts, otlptracegrpc.WithInsecure())
		} else {
			conf, err := client.NewConfig(params.FS, params.Config.TLS)
			if err != nil {
				return nil, err
			}
			opts = append(opts, otlptracegrpc.WithTLSCredentials(grpc.NewTLS(conf)))
		}
		if !strings.IsEmpty(params.Config.URL) {
			opts = append(opts, otlptracegrpc.WithEndpoint(params.Config.URL))
		}
		return otlptrace.NewUnstarted(otlptracegrpc.NewClient(opts...)), nil
	default:
		opts := []otlptracehttp.Option{otlptracehttp.WithHeaders(params.Config.Headers)}
		if !strings.IsEmpty(params.Config.URL) {
			opts = append(opts, otlptracehttp.WithEndpointURL(params.Config.URL))
		}
		return otlptrace.NewUnstarted(otlptracehttp.NewClient(opts...)), nil
	}
}

// IsEnabled reports whether this package has registered tracing as enabled.
func IsEnabled() bool {
	return enabled.Load()
}

func newSampler(cfg *SamplerConfig) (sdk.Sampler, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	if cfg.Ratio < 0 || cfg.Ratio > 1 {
		return nil, ErrInvalidSampler
	}

	switch cfg.Kind {
	case "always_on":
		return sdk.AlwaysSample(), nil
	case "always_off":
		return sdk.NeverSample(), nil
	case "ratio":
		return sdk.ParentBased(sdk.TraceIDRatioBased(cfg.Ratio)), nil
	default:
		return nil, ErrInvalidSampler
	}
}
