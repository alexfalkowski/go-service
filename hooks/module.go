package hooks

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires Standard Webhooks helpers into Fx/Dig.
//
// It provides constructors for:
//   - *Generator (via NewGenerator), which generates new secret values suitable for Standard Webhooks, and
//   - *standardwebhooks.Webhook (via NewHook), which constructs a webhook instance from configuration.
//
// Disabled behavior: if hooks configuration is disabled (nil *Config), NewHook returns (nil, nil) so
// downstream consumers can treat webhook verification/signing as optional.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewHook),
)
