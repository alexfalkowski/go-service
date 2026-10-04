package otlp

import (
	"github.com/alexfalkowski/go-service/v2/config/validate"
	"github.com/alexfalkowski/go-service/v2/time"
)

// newCadenceValidation declares the `otlp_cadence` validation rule owned by this package.
func newCadenceValidation() validate.ValidationResult {
	return validate.ValidationResult{Validation: validate.Validation{Tag: "otlp_cadence", Func: validateCadence}}
}

// newBatchConfigValidation declares the `otlp_batch_config` validation rule owned by this package.
func newBatchConfigValidation() validate.ValidationResult {
	return validate.ValidationResult{
		Validation: validate.Validation{Tag: "otlp_batch_config", Func: validateBatchConfig},
	}
}

func validateCadence(fl validate.FieldLevel) bool {
	if isOTLP(fl) {
		return ValidateCadence(time.Duration(fl.Field().Int())) == nil
	}

	return true
}

func validateBatchConfig(fl validate.FieldLevel) bool {
	if isOTLP(fl) {
		parent := fl.Parent()
		cfg := BatchConfig{
			MaxQueueSize:       int(parent.FieldByName("MaxQueueSize").Int()),
			MaxExportBatchSize: int(parent.FieldByName("MaxExportBatchSize").Int()),
		}
		return ValidateBatchConfig(cfg) == nil
	}

	return true
}

func isOTLP(fl validate.FieldLevel) bool {
	kind := fl.Parent().FieldByName("Kind")
	return kind.IsValid() && kind.String() == "otlp"
}
