package rand

import "github.com/alexfalkowski/go-signal/v2"

// Module wires the crypto/rand subsystem into [github.com/alexfalkowski/go-signal/v2].
//
// It provides constructors for:
//   - Reader (via NewReader), which returns a cryptographically secure random source (crypto/rand.Reader), and
//   - *[Generator] (via NewGenerator), which produces cryptographically secure random values derived from that Reader.
//
// This module does not require configuration.
var Module = signal.Module(
	signal.Constructor(NewReader),
	signal.Constructor(NewGenerator),
)
