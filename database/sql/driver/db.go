package driver

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/linxGnu/mssqlx"
)

// DBs wraps [mssqlx.DBs].
//
// It embeds the upstream type so callers keep the usual mssqlx query, exec,
// transaction, ping, and pool-configuration methods while go-service can attach
// repository-owned lifecycle cleanup such as OpenTelemetry metric registration
// cleanup.
type DBs struct {
	*mssqlx.DBs

	registrations []metrics.Registration
}

// Destroy unregisters repository-owned DB stats metrics and closes all database
// pools.
func (d *DBs) Destroy() error {
	if d == nil {
		return nil
	}

	regs := d.registrations
	d.registrations = nil

	errs := unregister(regs)
	if d.DBs != nil {
		errs = append(errs, d.DBs.Destroy()...)
	}

	return errors.Join(errs...)
}
