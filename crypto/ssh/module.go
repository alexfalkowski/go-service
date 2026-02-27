package ssh

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the SSH crypto subsystem into Fx/Dig.
//
// It provides constructors for:
//   - *Generator (via NewGenerator), which generates Ed25519 SSH key pairs,
//   - *Signer (via NewSigner), which signs messages when SSH config is enabled, and
//   - *Verifier (via NewVerifier), which verifies signatures when SSH config is enabled.
//
// Disabled behavior: if SSH configuration is disabled (nil *Config), NewSigner and NewVerifier
// return (nil, nil) so downstream consumers can treat SSH signing/verification as optional.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewSigner),
	di.Constructor(NewVerifier),
)
