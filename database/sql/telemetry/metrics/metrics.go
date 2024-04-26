package metrics

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(dbs *mssqlx.DBs, version version.Version, meter metric.Meter) {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
		attribute.Key("db_driver").String(dbs.DriverName()),
	)

	maxOpen := metrics.MustInt64ObservableGauge(meter, "sql_max_open_total", "Maximum number of open connections to the database.")
	open := metrics.MustInt64ObservableGauge(meter, "sql_open_total", "The number of established connections both in use and idle.")
	inUse := metrics.MustInt64ObservableGauge(meter, "sql_in_use_total", "The number of connections currently in use.")
	idle := metrics.MustInt64ObservableGauge(meter, "sql_idle_total", "The number of idle connections.")
	waited := metrics.MustInt64ObservableCounter(meter, "sql_waited_for_total", "The total number of connections waited for.")
	blocked := metrics.MustFloat64ObservableCounter(meter, "sql_blocked_seconds_total", "The total time blocked waiting for a new connection.")
	maxIdleClosed := metrics.MustInt64ObservableCounter(meter, "sql_closed_max_idle_total", "The total number of connections closed due to SetMaxIdleConns.")
	maxIdleTimeClosed := metrics.MustInt64ObservableCounter(meter, "sql_closed_max_lifetime_total", "The total number of connections closed due to SetConnMaxIdleTime.")
	maxLifetimeClosed := metrics.MustInt64ObservableCounter(meter, "sql_closed_max_idle_time_total", "The total number of connections closed due to SetConnMaxLifetime.")

	m := &ms{
		dbs: dbs, opts: opts,
		mo: maxOpen, o: open, iu: inUse,
		i: idle, w: waited, b: blocked,
		mic: maxIdleClosed, mitc: maxIdleTimeClosed, mlc: maxLifetimeClosed,
	}

	meter.RegisterCallback(m.callback, maxOpen, open, inUse, idle, waited, blocked, maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed)
}

type ms struct {
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

func (m *ms) callback(_ context.Context, o metric.Observer) error {
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

func (m *ms) collect(db *sqlx.DB, o metric.Observer, opts metric.MeasurementOption) {
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
