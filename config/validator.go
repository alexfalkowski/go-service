package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/otlp"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/go-playground/validator/v10"
)

// FieldLevel aliases the upstream validator field-level interface for custom config validation rules.
type FieldLevel = validator.FieldLevel

// NewValidator constructs a Validator backed by go-playground/validator.
//
// It enables required-struct validation ([validator.WithRequiredStructEnabled]), which causes validation
// tags like `required` to be applied to nested struct fields in a more strict/consistent way.
//
// It also registers repository-owned validation tags:
//   - `config_size`: accepts integer-like byte sizes between 0 and [bytes.MaxConfigSize].
//   - `duration_second_precision`: accepts positive durations that are exact multiples of one second.
//   - `otlp_cadence`: accepts OTLP exporter cadence settings for OTLP configurations.
//   - `otlp_batch_config`: accepts OTLP queue and export batch settings for OTLP configurations.
//
// This constructor is typically wired via [Module] and consumed by `NewConfig[T]` to validate
// decoded configuration before returning it to the caller.
func NewValidator() *Validator {
	validate := validator.New(validator.WithRequiredStructEnabled(), validator.WithTagNameFuncBlankOmit())
	runtime.Must(validate.RegisterValidation("config_size", validateConfigSize))
	runtime.Must(validate.RegisterValidation("duration_second_precision", validateDurationSecondPrecision))
	runtime.Must(validate.RegisterValidation("otlp_cadence", validateOTLPCadence))
	runtime.Must(validate.RegisterValidation("otlp_batch_config", validateOTLPBatchConfig))

	return &Validator{validate}
}

// Validator wraps a go-playground validator instance.
//
// It is used by `NewConfig[T]` to validate decoded configuration structs. You may use the embedded
// `*validator.Validate` directly to register custom validations or to validate values manually.
type Validator struct {
	*validator.Validate
}

func validateConfigSize(fl FieldLevel) bool {
	size := bytes.Size(fl.Field().Int())
	return size >= 0 && size <= bytes.MaxConfigSize
}

func validateDurationSecondPrecision(fl FieldLevel) bool {
	duration := time.Duration(fl.Field().Int())
	return duration > 0 && duration%time.Second == 0
}

func validateOTLPCadence(fl FieldLevel) bool {
	if isOTLP(fl) {
		return otlp.ValidateCadence(time.Duration(fl.Field().Int())) == nil
	}

	return true
}

func validateOTLPBatchConfig(fl FieldLevel) bool {
	if isOTLP(fl) {
		parent := fl.Parent()
		cfg := otlp.BatchConfig{
			MaxQueueSize:       int(parent.FieldByName("MaxQueueSize").Int()),
			MaxExportBatchSize: int(parent.FieldByName("MaxExportBatchSize").Int()),
		}
		return otlp.ValidateBatchConfig(cfg) == nil
	}

	return true
}

func isOTLP(fl FieldLevel) bool {
	kind := fl.Parent().FieldByName("Kind")
	return kind.IsValid() && kind.String() == "otlp"
}
