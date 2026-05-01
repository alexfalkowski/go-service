// Package access provides authorization (access control) helpers used by go-service.
//
// This package is focused on answering the question “is subject X allowed to do Y?”
// based on an authorization policy. It currently provides:
//
//   - Controller: an interface for permission checks.
//   - A Casbin-backed implementation of Controller.
//
// # Access checks
//
// Controller.HasAccess expects separate user, system, and action values. The
// implementation treats system as the authorization object and action as the
// operation being requested.
//
// # Casbin implementation
//
// NewController constructs a CasbinController backed by github.com/casbin/casbin/v2:
//
//   - The model is resolved with os.FS.ReadSource and parsed by Casbin.
//   - The policy is resolved with os.FS.ReadSource and loaded using Casbin's string adapter.
//
// Model and Policy support go-service source strings such as "env:" and "file:",
// plus literal content.
//
// # Enablement
//
// Enablement is modeled by presence: if *Config is nil, NewController returns (nil, nil).
// Callers should handle a nil Controller as “authorization disabled/unconfigured” and
// decide on an appropriate default behavior at a higher layer.
//
// Start with Config and NewController.
package access
