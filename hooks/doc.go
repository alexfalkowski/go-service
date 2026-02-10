// Package hooks provides Standard Webhooks helpers and wiring for go-service.
//
// This package contains:
//   - configuration for webhook secrets (loaded via os.FS.ReadSource), and
//   - constructors for creating Standard Webhooks hook instances and generating secrets.
//
// Transport integrations (for example `transport/http/hooks`) build on top of this package to sign and verify requests.
package hooks
