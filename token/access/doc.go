// Package access provides authorization (access control) helpers used by go-service.
//
// This package is focused on answering the question "is the verified request
// subject allowed to invoke this transport operation?" based on an authorization
// policy. It currently provides:
//
//   - Controller: an interface for permission checks.
//   - A Casbin-backed implementation of Controller.
//
// # Access checks
//
// [Controller.HasAccess] expects a request context populated by the transport
// metadata middleware. The Casbin request tuple is:
//
//   - sub: [github.com/alexfalkowski/go-service/v2/meta.UserID]
//   - obj: [github.com/alexfalkowski/go-service/v2/meta.TransportServiceMethod]
//   - act: "invoke"
//
// HTTP servers use the matched route pattern as the service-method when one is
// available, for example "http:GET /users/{id}". gRPC servers use the full RPC
// method, for example "grpc:/package.Service/Method".
//
// # Casbin implementation
//
// NewController constructs a CasbinController backed by [github.com/casbin/casbin/v2]:
//
//   - The model is resolved with [os.FS.ReadSource] and parsed by Casbin.
//   - The policy is resolved with [os.FS.ReadSource] and loaded using Casbin's string adapter.
//
// Security note: Casbin's string adapter ignores malformed policy-line errors.
// Policy content is treated as trusted deployment configuration here, and startup
// still fails for empty policy content or invalid model construction.
//
// Model and Policy support go-service source strings such as "env:NAME" and
// "file:<path>", plus literal content. File paths may be absolute or relative
// according to [os.FS.ReadSource].
//
// # Enablement
//
// Enablement is modeled by presence: if *[Config] is nil, NewController returns
// (nil, nil). Callers should handle a nil Controller as "authorization
// disabled/unconfigured" and decide on an appropriate default behavior at a
// higher layer.
//
// Start with Config and NewController.
package access
