package sql

import (
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires SQL database support into [go.uber.org/fx]/[go.uber.org/dig].
//
// This module is the aggregate SQL bundle used by higher-level application bundles such as
// [github.com/alexfalkowski/go-service/v2/module.Server] and [github.com/alexfalkowski/go-service/v2/module.Client]. It keeps SQL support opt-in via a single DI option while
// delegating backend-specific behavior to subpackages.
//
// At present it includes only PostgreSQL support via [github.com/alexfalkowski/go-service/v2/database/sql/pg.Module], but the module shape
// is intentionally extensible for additional SQL drivers in the future.
//
// Driver-specific constructors read DSNs/pool settings from config, open writer/reader pools, install
// OpenTelemetry driver instrumentation when tracing or metrics are enabled, register DB stats metrics
// when metrics are enabled, and attach lifecycle hooks that close pools on shutdown.
var Module = di.Module(
	pg.Module,
)
