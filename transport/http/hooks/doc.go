// Package hooks provides HTTP webhook middleware and wiring for go-service.
//
// This package adapts the Standard Webhooks library to go-service HTTP transports by providing:
//   - a `Webhook` wrapper that signs outbound requests and verifies inbound requests,
//   - server-side verification middleware via `NewHandler`, and
//   - client-side signing middleware via `NewRoundTripper`.
//
// Disabled behavior:
// When the underlying Standard Webhooks instance is nil, webhook support is treated as disabled and the
// transport helpers become pass-through no-ops instead of failing or panicking. This allows higher-level
// wiring such as CloudEvents integrations to keep webhook support optional.
//
// Start with `NewWebhook` to adapt a Standard Webhooks instance into transport middleware helpers.
package hooks
