// Package hooks provides shared Standard Webhooks construction helpers and wiring for go-service.
//
// This package integrates the Standard Webhooks Go library ([github.com/standard-webhooks/standard-webhooks/libraries/go])
// by providing:
//
//   - configuration for the named webhook secret set used by Standard Webhooks
//     signers/verifiers (see [Config]), loaded via the
//     go-service "source string" pattern (resolved by
//     [github.com/alexfalkowski/go-service/v2/os.FS.ReadSource]), and
//
//   - constructors for creating Standard Webhooks hooks (see [NewHook]) and generating new secrets (see [Generator]).
//
// # Secrets and source strings
//
// Each value in [Config.Secrets] is a "source string" and may be:
//
//   - "env:NAME" to read the secret from an environment variable,
//   - "file:/path/to/secret" to read the secret from a file, or
//   - any other value treated as the literal secret.
//
// The resolved secret bytes are passed to the Standard Webhooks library as a string. Empty resolved
// secrets are rejected by this package, while non-empty secrets must use a format accepted by the
// Standard Webhooks constructor. The [Generator] output is accepted directly as a literal secret. Keep
// this value private and avoid logging it.
//
// [Config.Key] selects the active entry in [Config.Secrets] used for signing.
// Verification accepts signatures from every configured secret
// because Standard Webhooks signs a message id, timestamp, and payload, but does
// not include a signing key id that could select one secret directly.
//
// # Downstream integrations
//
// Transport integrations (for example [github.com/alexfalkowski/go-service/v2/transport/http/hooks])
// build on top of this package to verify incoming webhook signatures and/or sign outbound webhook
// requests.
package hooks
