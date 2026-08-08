package otlp

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ErrInvalidCadence is returned when an OTLP export cadence is neither zero nor
// a positive whole number of seconds.
var ErrInvalidCadence = errors.New("otlp: invalid export cadence")

// ValidateCadence accepts zero, which leaves the OpenTelemetry SDK default in
// effect, and positive whole-second export cadences.
func ValidateCadence(cadence time.Duration) error {
	if cadence == 0 || cadence > 0 && cadence%time.Second == 0 {
		return nil
	}

	return ErrInvalidCadence
}
