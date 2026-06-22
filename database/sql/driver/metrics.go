package driver

import (
	"database/sql"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
)

func register(dbs []*sql.DB, role string, opts ...telemetry.Option) []metrics.Registration {
	regs := make([]metrics.Registration, 0, len(dbs))

	for i, db := range dbs {
		reg, err := telemetry.RegisterDBStatsMetrics(db, statsOptions(role, i, opts)...)
		runtime.Must(err)
		regs = append(regs, reg)
	}

	return regs
}

func statsOptions(role string, index int, opts []telemetry.Option) []telemetry.Option {
	options := make([]telemetry.Option, 0, len(opts)+1)
	options = append(options, opts...)

	name := role + "." + strconv.Itoa(index)
	options = append(options, telemetry.WithAttributes(attributes.DBClientConnectionPoolName(name)))

	return options
}

func unregister(regs []metrics.Registration) []error {
	errs := make([]error, 0, len(regs))
	for _, reg := range regs {
		errs = append(errs, reg.Unregister())
	}

	return errs
}

func options(name string, opts []telemetry.Option) []telemetry.Option {
	if len(opts) > 0 {
		return opts
	}

	return []telemetry.Option{telemetry.WithAttributes(attributes.DBSystemNameKey.String(name))}
}
