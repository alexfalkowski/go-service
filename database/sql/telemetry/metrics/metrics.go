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
	if dbs == nil {
		return
	}

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

	metrics := &Metrics{
		dbs: dbs, opts: opts, maxOpen: maxOpen, open: open,
		inUse: inUse, idle: idle, waited: waited, blocked: blocked,
		maxIdleClosed: maxIdleClosed, maxIdleTimeClosed: maxIdleTimeClosed, maxLifetimeClosed: maxLifetimeClosed,
	}

	meter.MustRegisterCallback(metrics.callback, maxOpen, open, inUse, idle, waited, blocked, maxIdleClosed, maxIdleTimeClosed, maxLifetimeClosed)
}

// Metrics for SQL.
type Metrics struct {
	dbs  *mssqlx.DBs
	opts metric.MeasurementOption

	maxOpen           metric.Int64ObservableGauge
	open              metric.Int64ObservableGauge
	inUse             metric.Int64ObservableGauge
	idle              metric.Int64ObservableGauge
	waited            metric.Int64ObservableCounter
	blocked           metric.Float64ObservableCounter
	maxIdleClosed     metric.Int64ObservableCounter
	maxIdleTimeClosed metric.Int64ObservableCounter
	maxLifetimeClosed metric.Int64ObservableCounter
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

	observer.ObserveInt64(m.maxOpen, int64(stats.MaxOpenConnections), m.opts, opts)
	observer.ObserveInt64(m.open, int64(stats.OpenConnections), m.opts, opts)
	observer.ObserveInt64(m.inUse, int64(stats.InUse), m.opts, opts)
	observer.ObserveInt64(m.idle, int64(stats.Idle), m.opts, opts)
	observer.ObserveInt64(m.waited, stats.WaitCount, m.opts, opts)
	observer.ObserveFloat64(m.blocked, stats.WaitDuration.Seconds(), m.opts, opts)
	observer.ObserveInt64(m.maxIdleClosed, stats.MaxIdleClosed, m.opts, opts)
	observer.ObserveInt64(m.maxIdleTimeClosed, stats.MaxIdleTimeClosed, m.opts, opts)
	observer.ObserveInt64(m.maxLifetimeClosed, stats.MaxLifetimeClosed, m.opts, opts)
}
