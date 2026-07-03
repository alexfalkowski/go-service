package telemetry

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/errors"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	propagatorB3           = "b3"
	propagatorB3Multi      = "b3multi"
	propagatorBaggage      = "baggage"
	propagatorNone         = "none"
	propagatorTraceContext = "tracecontext"
)

var defaultPropagators = []string{propagatorTraceContext, propagatorBaggage}

// ErrInvalidPropagator is returned when propagation config contains an unsupported name.
var ErrInvalidPropagator = errors.New("telemetry: invalid propagator")

// PropagationConfig configures OpenTelemetry context propagation.
type PropagationConfig struct {
	// Formats lists the propagators used to extract and inject context.
	//
	// When empty, propagation defaults to W3C Trace Context and Baggage. Use
	// "none" as the sole value to disable propagation.
	Formats []string `yaml:"formats,omitempty" json:"formats,omitempty" toml:"formats,omitempty" validate:"omitempty,dive,oneof=tracecontext baggage b3 b3multi none"`
}

// RegisterPropagation configures the global OpenTelemetry TextMapPropagator.
func RegisterPropagation(propagator propagation.TextMapPropagator) {
	otel.SetTextMapPropagator(propagator)
}

// NewPropagator constructs an OpenTelemetry propagator from cfg.
func NewPropagator(cfg *PropagationConfig) (propagation.TextMapPropagator, error) {
	propagators, err := newPropagators(propagationFormats(cfg))
	if err != nil {
		return nil, err
	}

	return propagation.NewCompositeTextMapPropagator(propagators...), nil
}

func propagationFormats(cfg *PropagationConfig) []string {
	if cfg != nil && len(cfg.Formats) > 0 {
		return cfg.Formats
	}
	return defaultPropagators
}

func newPropagators(formats []string) ([]propagation.TextMapPropagator, error) {
	if len(formats) == 1 && formats[0] == propagatorNone {
		return nil, nil
	}
	if slices.Contains(formats, propagatorNone) {
		return nil, ErrInvalidPropagator
	}

	propagators := make([]propagation.TextMapPropagator, 0, len(formats))
	for _, format := range formats {
		switch format {
		case propagatorTraceContext:
			propagators = append(propagators, propagation.TraceContext{})
		case propagatorBaggage:
			propagators = append(propagators, propagation.Baggage{})
		case propagatorB3:
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)))
		case propagatorB3Multi:
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
		default:
			return nil, ErrInvalidPropagator
		}
	}

	return propagators, nil
}
