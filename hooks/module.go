package hooks

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires shared Standard Webhooks helpers into [go.uber.org/fx]/[go.uber.org/dig].
//
// It provides constructors for:
//   - *[Generator] (via [NewGenerator]), which generates new secret values suitable for Standard Webhooks, and
//   - *[Hook] (via [NewHook]), which constructs a signer/verifier from configuration.
//
// Disabled behavior: if hooks configuration is disabled (nil *[Config]), [NewHook] returns (nil, nil) so
// downstream consumers can treat webhook verification/signing as optional.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewHook),
)
