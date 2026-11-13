package ed25519

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the Ed25519 subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//   - *[Generator] (via NewGenerator), which generates Ed25519 key pairs,
//   - *[Signer] (via NewSigner), which signs messages when Ed25519 config is enabled, and
//   - *[Verifier] (via NewVerifier), which verifies signatures when Ed25519 config is enabled.
//
// Disabled behavior: if Ed25519 configuration is disabled (nil *[Config]), NewSigner and NewVerifier
// return (nil, nil) so downstream consumers can treat signing/verification as optional.
var Module = signal.Module(
	signal.Constructor(NewGenerator),
	signal.Constructor(NewSigner),
	signal.Constructor(NewVerifier),
)
