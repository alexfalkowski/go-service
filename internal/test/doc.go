// Package test provides shared fixtures, configuration builders, and in-memory
// wiring helpers used by integration-style tests across the repository.
//
// The package mirrors production composition closely enough to exercise HTTP,
// gRPC, telemetry, health, cache, database, token, and debug flows without
// forcing each test to rebuild that setup from scratch.
package test
