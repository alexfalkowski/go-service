// Package hooks provides HTTP webhook middleware and wiring for go-service.
//
// This package adapts the Standard Webhooks library to go-service HTTP transports by providing:
//   - a [Webhook] wrapper that signs outbound requests and verifies inbound requests,
//   - server-side verification middleware via [NewHandler], and
//   - client-side signing middleware via [NewRoundTripper].
//
// Disabled behavior:
// When the underlying Standard Webhooks instance is nil, webhook support is treated as disabled and the
// transport helpers become pass-through no-ops instead of failing or panicking. This allows higher-level
// wiring such as CloudEvents integrations to keep webhook support optional.
//
// Body limits:
// The webhook helpers buffer request bodies so signatures can be computed and
// verified while leaving req.Body readable for downstream handlers. Under the
// supported transport wiring, inbound webhook requests are already capped by the
// HTTP server request-body limiter configured from MaxReceiveSize before mux
// handlers run. Outbound requests are created by callers and are expected to be
// bounded at the request construction boundary.
//
// Replay protection:
// Verification covers the Standard Webhooks signature and timestamp checks, but
// this package does not keep replay state or reject a previously seen
// Webhook-Id. Receivers that perform non-idempotent work must deduplicate or
// process idempotently using the Webhook-Id or event identifier, preferably with
// durable storage shared by all receiver instances.
//
// Start with [NewWebhook] to adapt a Standard Webhooks instance into transport middleware helpers.
package hooks
