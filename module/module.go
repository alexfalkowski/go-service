package module

import (
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/compress"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

var (
	// Library provides a baseline Fx module intended for reuse by both servers and clients.
	//
	// It wires common, transport-agnostic dependencies that many subsystems build upon:
	//   - env.Module (service identity values like name/version/id, user agent, etc.)
	//   - compress.Module (compression registry and default codecs)
	//   - encoding.Module (encoding registry and default encoders)
	//   - crypto.Module (crypto primitives and helpers)
	//   - time.Module (time providers/utilities)
	//   - sync.Module (shared buffer pool wiring)
	//   - id.Module (ID generator implementations and selection)
	//
	// Library does not wire transports, servers, or request handling; it is intended to be a common
	// foundation that both the Server and Client bundles build upon.
	Library = di.Module(
		env.Module,
		compress.Module,
		encoding.Module,
		crypto.Module,
		time.Module,
		sync.Module,
		id.Module,
	)

	// Server provides the standard Fx module composition for a go-service server.
	//
	// It builds on Library and adds server-oriented wiring commonly needed by services:
	//   - debug.Module (debug server + diagnostic endpoints)
	//   - config.Module (config decoding + validation + common sub-config projections)
	//   - telemetry.Module (logging/tracing/metrics wiring)
	//   - cache.Module (cache drivers, cache facade, and package-level cache registration)
	//   - feature.Module (OpenFeature client + optional provider registration)
	//   - sql.Module (SQL database wiring; currently PostgreSQL)
	//   - limiter.Module (rate limiter key map wiring; transport modules typically construct limiters)
	//   - transport.Module (HTTP/gRPC transports and related client/server integrations)
	//   - health.Module (health server wiring)
	//
	// Many of these subsystems are optional and are enabled/disabled by configuration (often via nil pointer
	// sub-configs). This bundle wires constructors/registrations; runtime behavior depends on the config
	// supplied to the graph.
	//
	// This is the primary entrypoint for server applications built from `go-service-template`.
	Server = di.Module(
		Library,
		debug.Module,
		config.Module,
		telemetry.Module,
		cache.Module,
		feature.Module,
		sql.Module,
		limiter.Module,
		transport.Module,
		health.Module,
	)

	// Client provides the standard Fx module composition for a go-service client.
	//
	// It builds on Library and adds client-oriented wiring commonly needed by client processes
	// and batch jobs:
	//   - config.Module (config decoding + validation + common sub-config projections)
	//   - telemetry.Module (logging/tracing/metrics wiring)
	//   - cache.Module (cache drivers/facade; optional by config)
	//   - feature.Module (OpenFeature client + optional provider registration)
	//   - hooks.Module (Standard Webhooks helpers; used by HTTP hook integrations)
	//   - sql.Module (SQL database wiring; currently PostgreSQL)
	//   - limiter.Module (rate limiter key map wiring)
	//
	// Unlike Server, Client does not wire debug endpoints, transports, or a health server by default.
	// Those can be added explicitly by composing additional modules on top of Client if needed.
	//
	// This is the primary entrypoint for client-style applications built from `go-service-template`.
	Client = di.Module(
		Library,
		config.Module,
		telemetry.Module,
		cache.Module,
		feature.Module,
		hooks.Module,
		sql.Module,
		limiter.Module,
	)
)
