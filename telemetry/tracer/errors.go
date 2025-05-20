package tracer

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Error for tracer.
func Error(err error, span trace.Span) {
	if err == nil {
		return
	}

	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

func prefix(err error) error {
	return errors.Prefix("metrics", err)
}
