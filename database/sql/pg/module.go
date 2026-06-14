package pg

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires PostgreSQL ([github.com/alexfalkowski/go-service/v2/database/sql]) support into [go.uber.org/fx]/[go.uber.org/dig].
//
// It registers the [github.com/jackc/pgx/v5/stdlib] driver under the name "pg" (with OpenTelemetry instrumentation
// when tracing or metrics are enabled) and provides a constructor that opens master/slave connection pools.
//
// Provided components:
//   - Registration: [Register]
//   - Constructor: [Open]
//
// Disabled behavior: [Open] returns (nil, nil) when PostgreSQL configuration is disabled.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(Open),
)
