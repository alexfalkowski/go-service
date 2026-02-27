package limiter

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires limiter key derivation helpers into Fx/Dig.
//
// It provides a constructor for the default `KeyMap` via `NewKeyMap`. The returned map contains the
// built-in key kinds supported by this package (e.g. "user-agent", "ip", and "token").
//
// # Extending key kinds
//
// If you need additional key kinds, you can provide your own `KeyMap` instead of using this module,
// or decorate/replace the provided map in your application wiring and add new entries before it is
// used to construct limiters.
//
// This module intentionally does not construct a `*Limiter`; limiter construction is transport- and
// feature-specific and typically requires configuration (see `NewLimiter` and `Config`).
var Module = di.Module(
	di.Constructor(NewKeyMap),
)
