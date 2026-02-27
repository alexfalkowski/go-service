package rand

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires the crypto/rand subsystem into Fx/Dig.
//
// It provides constructors for:
//   - Reader (via NewReader), which returns a cryptographically secure random source (crypto/rand.Reader), and
//   - *Generator (via NewGenerator), which produces cryptographically secure random values derived from that Reader.
//
// This module does not require configuration.
var Module = di.Module(
	di.Constructor(NewReader),
	di.Constructor(NewGenerator),
)
