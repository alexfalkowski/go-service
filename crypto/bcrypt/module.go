package bcrypt

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the bcrypt subsystem into Fx/Dig.
//
// It provides a constructor for `*Signer` via `NewSigner`, which can be used for password hashing
// and verification using bcrypt with the default cost.
//
// This module does not require configuration.
var Module = di.Module(
	di.Constructor(NewSigner),
)
