package time

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the network time provider constructor and this package's `duration_second_precision`
// validation rule into [github.com/alexfalkowski/go-signal/v2].
//
// Including this module in a go-signal application provides a constructor for Network via
// NewNetwork.
//
// NewNetwork uses *[Config] to decide whether to enable network time and which provider
// to construct (for example "ntp" or "nts"). When network time is disabled (nil config),
// the constructor returns (nil, nil). A non-nil config with an empty or unknown Kind
// is treated as enabled but invalid, so the constructor returns [ErrNotFound].
//
// This module does not force the application to use network time; it only makes the
// provider available for optional injection.
//
// It also contributes a [github.com/alexfalkowski/go-service/v3/config/validate.Validation] to the
// shared "validations" group consumed by
// [github.com/alexfalkowski/go-service/v3/config/validate.NewValidator].
var Module = signal.Module(
	signal.Constructor(NewNetwork),
	signal.Constructor(newValidation),
)
