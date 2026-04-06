// Package driver provides low-level SQL driver registration and connection helpers for go-service.
//
// This package is the bridge between driver-specific packages (for example `database/sql/pg`)
// and the shared SQL wiring used by the higher-level module graph.
//
// Its main responsibilities are:
//   - `Register`, which wraps a concrete `database/sql/driver.Driver` with OpenTelemetry
//     instrumentation and installs it in the global `database/sql` driver registry.
//   - `Open`, which resolves DSNs from go-service source strings, opens master/slave pools,
//     applies pool settings, registers DB stats metrics, and closes those pools on lifecycle stop.
//
// Most applications use this package indirectly through a driver package such as
// `database/sql/pg` rather than calling it directly.
package driver
