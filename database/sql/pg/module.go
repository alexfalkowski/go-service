package pg

import "github.com/alexfalkowski/go-signal/v2"

// Module wires PostgreSQL ([github.com/alexfalkowski/go-service/v3/database/sql]) support into [github.com/alexfalkowski/go-signal/v2].
//
// It registers the [github.com/jackc/pgx/v5/stdlib] driver under the name "pg" (with OpenTelemetry instrumentation
// when tracing or metrics are enabled) and provides a constructor that opens writer/reader connection pools.
//
// Provided components:
//   - Registration: [Register]
//   - Constructor: [Open]
//
// Disabled behavior: [Open] returns (nil, nil) when PostgreSQL configuration is disabled.
var Module = signal.Module(
	signal.Register(Register),
	signal.Constructor(Open),
)
