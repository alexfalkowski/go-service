package time

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the network time provider constructor into [go.uber.org/fx].
//
// Including this module in an Fx application provides a constructor for Network via
// NewNetwork.
//
// NewNetwork uses *[Config] to decide whether to enable network time and which provider
// to construct (for example "ntp" or "nts"). When network time is disabled (nil config),
// the constructor returns (nil, nil). A non-nil config with an empty or unknown Kind
// is treated as enabled but invalid, so the constructor returns [ErrNotFound].
//
// This module does not force the application to use network time; it only makes the
// provider available for optional injection.
var Module = di.Module(
	di.Constructor(NewNetwork),
)
