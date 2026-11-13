package bcrypt

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the bcrypt subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides a constructor for [*Signer] via [NewSigner], which can be used for password hashing
// and verification using bcrypt with the default cost.
//
// This module does not require configuration.
var Module = signal.Module(
	signal.Constructor(NewSigner),
)
