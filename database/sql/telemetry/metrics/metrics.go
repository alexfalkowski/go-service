package metrics

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
//
//nolint:funlen,cyclop
func Register(dbs *mssqlx.DBs, version version.Version, meter metric.Meter) error {
	if dbs == nil {
		return nil
	}

	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
		attribute.Key("db_driver").String(dbs.DriverName()),
	)

	maxOpen, err := meter.Int64ObservableGauge("sql_max_open_total", metric.WithDescription("Maximum number of open connections to the database."))
	if err != nil {
		return err
	}

	open, err := meter.Int64ObservableGauge("sql_open_total", metric.WithDescription("The number of established connections both in use and idle."))
	if err != nil {
		return err
	}

	inUse, err := meter.Int64ObservableGauge("sql_in_use_total", metric.WithDescription("The number of connections currently in use."))
	if err != nil {
		return err
	}

	idle, err := meter.Int64ObservableGauge("sql_idle_total", metric.WithDescription("The number of idle connections."))
	if err != nil {
		return err
	}

	waited, err := meter.Int64ObservableCounter("sql_waited_for_total", metric.WithDescription("The total number of connections waited for."))
	if err != nil {
		return err
	}

	blocked, err := meter.Float64ObservableCounter("sql_blocked_seconds_total", metric.WithDescription("The total time blocked waiting for a new connection."))
	if err != nil {
		return err
	}

	maxIdleClosed, err := meter.Int64ObservableCounter("sql_closed_max_idle_total", metric.WithDescription("The total number of connections closed due to SetMaxIdleConns."))
	if err != nil {
		return err
	}

	maxIdleTimeClosed, err := meter.Int64ObservableCounter("sql_closed_max_lifetime_total",
		metric.WithDescription("The total number of connections closed due to SetConnMaxIdleTime."))
	if err != nil {
		return err
	}

	maxLifetimeClosed, err := meter.Int64ObservableCounter("sql_closed_max_idle_time_total",
		metric.WithDescription("The total number of connections closed due to SetConnMaxLifetime."))
	if err != nil {
		return err
	}

	m := &metrics{
		dbs: dbs, opts: opts,
		mo: maxOpen, o: open, iu: inUse,
		i: idle, w: waited, b: blocked,
		mic: maxIdleClosed, mitc: maxIdleTimeClosed, mlc: maxLifetimeClosed,
	}

	meter.RegisterCallback(m.callback, maxOpen, open, inUse, idle, waited, blocked, maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed)

	return nil
}

type metrics struct {
	dbs  *mssqlx.DBs
	opts metric.MeasurementOption

	mo   metric.Int64ObservableGauge
	o    metric.Int64ObservableGauge
	iu   metric.Int64ObservableGauge
	i    metric.Int64ObservableGauge
	w    metric.Int64ObservableCounter
	b    metric.Float64ObservableCounter
	mic  metric.Int64ObservableCounter
	mitc metric.Int64ObservableCounter
	mlc  metric.Int64ObservableCounter
}

func (m *metrics) callback(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		m.collect(ma, o, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		m.collect(s, o, opts)
	}

	return nil
}

func (m *metrics) collect(db *sqlx.DB, o metric.Observer, opts metric.MeasurementOption) {
	stats := db.Stats()

	o.ObserveInt64(m.mo, int64(stats.MaxOpenConnections), m.opts, opts)
	o.ObserveInt64(m.o, int64(stats.OpenConnections), m.opts, opts)
	o.ObserveInt64(m.iu, int64(stats.InUse), m.opts, opts)
	o.ObserveInt64(m.i, int64(stats.Idle), m.opts, opts)
	o.ObserveInt64(m.w, stats.WaitCount, m.opts, opts)
	o.ObserveFloat64(m.b, stats.WaitDuration.Seconds(), m.opts, opts)
	o.ObserveInt64(m.mic, stats.MaxIdleClosed, m.opts, opts)
	o.ObserveInt64(m.mitc, stats.MaxIdleTimeClosed, m.opts, opts)
	o.ObserveInt64(m.mlc, stats.MaxLifetimeClosed, m.opts, opts)
}
