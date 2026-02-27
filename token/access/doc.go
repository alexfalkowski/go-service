// Package access provides authorization (access control) helpers used by go-service.
//
// This package is focused on answering the question “is subject X allowed to do Y?”
// based on an authorization policy. It currently provides:
//
//   - Controller: an interface for permission checks.
//   - A Casbin-backed implementation of Controller.
//   - A default RBAC model definition (ModelConfig) suitable for simple subject/object/action checks.
//
// # Permission strings
//
// Controller.HasAccess expects the permission string to be in the form:
//
//	<system>:<action>
//
// Example: "service:read".
//
// The implementation splits this string on the first ":" and treats the parts as the
// object (system) and action components in the authorization request.
//
// # Casbin implementation
//
// NewController constructs a CasbinController backed by github.com/casbin/casbin/v2:
//
//   - The model is created from the embedded ModelConfig.
//   - The policy is loaded using Casbin’s file adapter.
//
// Note: despite the adapter name, the policy string is passed directly to the adapter.
// Ensure the configured policy value matches what the underlying adapter expects for
// your deployment (for example a path vs. a literal policy payload).
//
// # Enablement
//
// Enablement is modeled by presence: if *Config is nil, NewController returns (nil, nil).
// Callers should handle a nil Controller as “authorization disabled/unconfigured” and
// decide on an appropriate default behavior at a higher layer.
//
// Start with Config and NewController.
package access
