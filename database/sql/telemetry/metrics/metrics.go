package metrics

import (
	"context"
	"strconv"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
func Register(dbs *mssqlx.DBs, meter *metrics.Meter) {
	opts := metric.WithAttributes(attribute.Key("db_driver").String(dbs.DriverName()))
	maxOpen := meter.MustInt64ObservableGauge("sql_max_open_total", "Maximum number of open connections to the database.")
	open := meter.MustInt64ObservableGauge("sql_open_total", "The number of established connections both in use and idle.")
	inUse := meter.MustInt64ObservableGauge("sql_in_use_total", "The number of connections currently in use.")
	idle := meter.MustInt64ObservableGauge("sql_idle_total", "The number of idle connections.")
	waited := meter.MustInt64ObservableCounter("sql_waited_for_total", "The total number of connections waited for.")
	blocked := meter.MustFloat64ObservableCounter("sql_blocked_seconds_total", "The total time blocked waiting for a new connection.")
	maxIdleClosed := meter.MustInt64ObservableCounter("sql_closed_max_idle_total", "The total number of connections closed due to SetMaxIdleConns.")
	maxIdleTimeClosed := meter.MustInt64ObservableCounter("sql_closed_max_lifetime_total", "The total number of connections closed due to SetConnMaxIdleTime.")
	maxLifetimeClosed := meter.MustInt64ObservableCounter("sql_closed_max_idle_time_total", "The total number of connections closed due to SetConnMaxLifetime.")

	mts := &Metrics{
		dbs: dbs, opts: opts,
		mo: maxOpen, o: open, iu: inUse,
		i: idle, w: waited, b: blocked,
		mic: maxIdleClosed, mitc: maxIdleTimeClosed, mlc: maxLifetimeClosed,
	}

	meter.MustRegisterCallback(mts.callback, maxOpen, open, inUse, idle, waited, blocked, maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed)
}

// Metrics for SQL.
type Metrics struct {
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

func (m *Metrics) callback(_ context.Context, observer metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		opts := metric.WithAttributes(
			attribute.Key("db_name").String("master_" + strconv.Itoa(i)),
		)

		m.collect(ma, observer, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		opts := metric.WithAttributes(
			attribute.Key("db_name").String("slave_" + strconv.Itoa(i)),
		)

		m.collect(s, observer, opts)
	}

	return nil
}

func (m *Metrics) collect(db *sqlx.DB, observer metric.Observer, opts metric.MeasurementOption) {
	stats := db.Stats()

	observer.ObserveInt64(m.mo, int64(stats.MaxOpenConnections), m.opts, opts)
	observer.ObserveInt64(m.o, int64(stats.OpenConnections), m.opts, opts)
	observer.ObserveInt64(m.iu, int64(stats.InUse), m.opts, opts)
	observer.ObserveInt64(m.i, int64(stats.Idle), m.opts, opts)
	observer.ObserveInt64(m.w, stats.WaitCount, m.opts, opts)
	observer.ObserveFloat64(m.b, stats.WaitDuration.Seconds(), m.opts, opts)
	observer.ObserveInt64(m.mic, stats.MaxIdleClosed, m.opts, opts)
	observer.ObserveInt64(m.mitc, stats.MaxIdleTimeClosed, m.opts, opts)
	observer.ObserveInt64(m.mlc, stats.MaxLifetimeClosed, m.opts, opts)
}
