// Package metrics provides HTTP metrics endpoint wiring for go-service.
//
// This package exposes a Prometheus HTTP handler under the service's HTTP mux when metrics are enabled
// and the configured kind is Prometheus.
//
// Start with `Register`.
package metrics
