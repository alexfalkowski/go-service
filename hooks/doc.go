// Package hooks provides Standard Webhooks helpers and wiring for go-service.
//
// This package integrates the Standard Webhooks Go library (github.com/standard-webhooks/standard-webhooks/libraries/go)
// by providing:
//
//   - configuration for webhook signing/verification secrets (see Config), loaded via the go-service "source string"
//     pattern (resolved by os.FS.ReadSource), and
//
//   - constructors for creating Standard Webhooks webhook instances (see NewHook) and generating new secrets (see Generator).
//
// # Secrets and source strings
//
// The secret configured in Config.Secret is a "source string" and may be:
//
//   - "env:NAME" to read the secret from an environment variable,
//   - "file:/path/to/secret" to read the secret from a file, or
//   - any other value treated as the literal secret.
//
// The resolved secret bytes are passed to the Standard Webhooks library as a string. Keep this value private and avoid
// logging it.
//
// # Downstream integrations
//
// Transport integrations (for example `transport/http/hooks`) build on top of this package to verify incoming webhook
// signatures and/or sign outbound webhook requests.
package hooks
