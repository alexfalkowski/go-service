package time

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the network time provider constructor into Fx.
//
// Including this module in an Fx application provides a constructor for Network via
// NewNetwork.
//
// NewNetwork uses *Config to decide whether to enable network time and which provider
// to construct (for example "ntp" or "nts"). When network time is disabled (nil config),
// the constructor returns (nil, nil).
//
// This module does not force the application to use network time; it only makes the
// provider available for optional injection.
var Module = di.Module(
	di.Constructor(NewNetwork),
)
