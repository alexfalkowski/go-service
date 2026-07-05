// Package driver provides low-level SQL driver registration and connection helpers for go-service.
//
// This package is the bridge between driver-specific packages (for example [github.com/alexfalkowski/go-service/v2/database/sql/pg])
// and the shared SQL wiring used by the higher-level module graph.
//
// Its main responsibilities are:
//   - [Register], which wraps a concrete [database/sql/driver.Driver] with OpenTelemetry
//     instrumentation and installs it in the global [database/sql] driver registry.
//   - [Open], which resolves DSNs from go-service source strings, creates writer/reader pools,
//     applies each role's pool settings, registers DB stats metrics, and closes those pools on lifecycle stop.
//
// Pool creation follows the [database/sql.Open] contract: it may not establish a
// network connection immediately. Use [github.com/alexfalkowski/go-service/v2/database/sql/driver.DBs.Ping],
// [github.com/alexfalkowski/go-service/v2/database/sql/driver.DBs.PingWriter],
// [github.com/alexfalkowski/go-service/v2/database/sql/driver.DBs.PingReader], or a
// health checker when startup/readiness must verify database reachability.
//
// Most applications use this package indirectly through a driver package such as
// [github.com/alexfalkowski/go-service/v2/database/sql/pg] rather than calling it directly.
package driver
