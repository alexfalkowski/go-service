// Package crypto provides cryptographic configuration and DI wiring used by go-service.
//
// This package is primarily an entrypoint for wiring multiple cryptographic subpackages into a single
// module (see `Module`) and for composing their configuration into `crypto.Config`.
//
// # Scope
//
// Most concrete cryptographic functionality lives in subpackages (for example `crypto/aes`,
// `crypto/ed25519`, `crypto/hmac`, `crypto/rsa`, `crypto/ssh`, `crypto/tls`, `crypto/pem`, and
// `crypto/rand`). Prefer importing those packages directly when you need specific primitives.
//
// The root package intentionally stays small:
//   - it defines the top-level `Config` used to enable/configure crypto subsystems, and
//   - it exports an Fx/Dig module that wires the supported crypto implementations.
//
// # Configuration conventions
//
// Sub-config fields on `crypto.Config` are pointers and are treated as optional. A nil sub-config is
// generally interpreted as "disabled" by the corresponding subsystem (see each subpackage's `IsEnabled`
// convention where applicable).
package crypto
