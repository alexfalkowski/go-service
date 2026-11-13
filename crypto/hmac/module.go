package hmac

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the HMAC subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//   - *[Generator] (via NewGenerator), which generates secret material suitable for HMAC keys, and
//   - *[Signer] (via NewSigner), which signs and verifies messages using HMAC-SHA-512 when HMAC config is enabled.
//
// Disabled behavior: if HMAC configuration is disabled (nil *[Config]), NewSigner returns (nil, nil) so
// downstream consumers can treat HMAC signing/verification as optional.
var Module = signal.Module(
	signal.Constructor(NewGenerator),
	signal.Constructor(NewSigner),
)
