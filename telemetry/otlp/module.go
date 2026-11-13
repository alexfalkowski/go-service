package otlp

import "github.com/alexfalkowski/go-signal/v2"

// Module wires this package's `otlp_cadence` and `otlp_batch_config` validation rules into
// [github.com/alexfalkowski/go-signal/v2].
//
// Including this module contributes two [github.com/alexfalkowski/go-service/v3/config/validate.Validation]
// values to the shared "validations" group consumed by
// [github.com/alexfalkowski/go-service/v3/config/validate.NewValidator].
var Module = signal.Module(
	signal.Constructor(newCadenceValidation),
	signal.Constructor(newBatchConfigValidation),
)
