package driver

import (
	"database/sql"

	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// Register registers a [database/sql] driver under name.
//
// This function registers the driver with the global [database/sql] driver registry. It is therefore intended
// to be called during process initialization (for example from an init hook or DI registration).
//
// Telemetry:
//   - The driver is wrapped using [telemetry.WrapDriver] when tracing or metrics are enabled.
//   - If opts is empty, the DB system name attribute is set to the provided name ([attributes.DBSystemNameKey]).
//
// Errors:
//   - If the underlying [sql.Register] panics (for example, due to registering the same name more than once),
//     Register converts that panic into an error and returns it.
func Register(name string, driver Driver, opts ...telemetry.Option) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	if metrics.IsEnabled() || tracer.IsEnabled() {
		driver = telemetry.WrapDriver(driver, options(name, opts)...)
	}

	sql.Register(name, driver)

	return err
}
