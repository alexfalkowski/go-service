package pg

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires PostgreSQL (`database/sql`) support into Fx/Dig.
//
// It registers the pgx stdlib driver under the name "pg" (with OpenTelemetry instrumentation)
// and provides a constructor that opens master/slave connection pools.
//
// Provided components:
//   - Registration: `Register`
//   - Constructor: `Open`
//
// Disabled behavior: `Open` returns (nil, nil) when PostgreSQL configuration is disabled.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(Open),
)
