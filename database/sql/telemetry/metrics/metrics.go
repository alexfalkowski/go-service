package metrics

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Register for metrics.
//
//nolint:funlen
func Register(dbs *mssqlx.DBs, version version.Version, meter metric.Meter) error {
	opts := metric.WithAttributes(
		attribute.Key("name").String(os.ExecutableName()),
		attribute.Key("version").String(string(version)),
		attribute.Key("db_driver").String(dbs.DriverName()),
	)

	maxOpen, err := meter.Float64ObservableGauge("sql_max_open_total", metric.WithDescription("Maximum number of open connections to the database."))
	if err != nil {
		return err
	}

	open, err := meter.Float64ObservableGauge("sql_open_total", metric.WithDescription("The number of established connections both in use and idle."))
	if err != nil {
		return err
	}

	inUse, err := meter.Float64ObservableGauge("sql_in_use_total", metric.WithDescription("The number of connections currently in use."))
	if err != nil {
		return err
	}

	idle, err := meter.Float64ObservableGauge("sql_idle_total", metric.WithDescription("The number of idle connections."))
	if err != nil {
		return err
	}

	waited, err := meter.Float64ObservableCounter("sql_waited_for_total", metric.WithDescription("The total number of connections waited for."))
	if err != nil {
		return err
	}

	blocked, err := meter.Float64ObservableCounter("sql_blocked_seconds_total", metric.WithDescription("The total time blocked waiting for a new connection."))
	if err != nil {
		return err
	}

	maxIdleClosed, err := meter.Float64ObservableCounter("sql_closed_max_idle_total", metric.WithDescription("The total number of connections closed due to SetMaxIdleConns."))
	if err != nil {
		return err
	}

	maxIdleTimeClosed, err := meter.Float64ObservableCounter("sql_closed_max_lifetime_total",
		metric.WithDescription("The total number of connections closed due to SetConnMaxIdleTime."))
	if err != nil {
		return err
	}

	maxLifetimeClosed, err := meter.Float64ObservableCounter("sql_closed_max_idle_time_total",
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

	meter.RegisterCallback(m.maxOpen, maxOpen)
	meter.RegisterCallback(m.open, open)
	meter.RegisterCallback(m.inUse, inUse)
	meter.RegisterCallback(m.idle, idle)
	meter.RegisterCallback(m.waited, waited)
	meter.RegisterCallback(m.blocked, blocked)
	meter.RegisterCallback(m.maxIdleClosed, maxIdleClosed)
	meter.RegisterCallback(m.maxIdleTimeClosed, maxIdleTimeClosed)
	meter.RegisterCallback(m.maxLifetimeClosed, maxLifetimeClosed)

	return nil
}

type metrics struct {
	dbs  *mssqlx.DBs
	opts metric.MeasurementOption

	mo   metric.Float64ObservableGauge
	o    metric.Float64ObservableGauge
	iu   metric.Float64ObservableGauge
	i    metric.Float64ObservableGauge
	w    metric.Float64ObservableCounter
	b    metric.Float64ObservableCounter
	mic  metric.Float64ObservableCounter
	mitc metric.Float64ObservableCounter
	mlc  metric.Float64ObservableCounter
}

func (m *metrics) maxOpen(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.mo, float64(stats.MaxOpenConnections), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.mo, float64(stats.MaxOpenConnections), m.opts, opts)
	}

	return nil
}

func (m *metrics) open(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.o, float64(stats.OpenConnections), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.o, float64(stats.OpenConnections), m.opts, opts)
	}

	return nil
}

func (m *metrics) inUse(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.iu, float64(stats.InUse), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.iu, float64(stats.InUse), m.opts, opts)
	}

	return nil
}

func (m *metrics) idle(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.i, float64(stats.Idle), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.i, float64(stats.Idle), m.opts, opts)
	}

	return nil
}

func (m *metrics) waited(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.WaitCount), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.WaitCount), m.opts, opts)
	}

	return nil
}

func (m *metrics) blocked(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.w, stats.WaitDuration.Seconds(), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.w, stats.WaitDuration.Seconds(), m.opts, opts)
	}

	return nil
}

func (m *metrics) maxIdleClosed(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxIdleClosed), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxIdleClosed), m.opts, opts)
	}

	return nil
}

func (m *metrics) maxIdleTimeClosed(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxIdleTimeClosed), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxIdleTimeClosed), m.opts, opts)
	}

	return nil
}

func (m *metrics) maxLifetimeClosed(_ context.Context, o metric.Observer) error {
	ms, _ := m.dbs.GetAllMasters()
	for i, ma := range ms {
		stats := ma.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("master_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxLifetimeClosed), m.opts, opts)
	}

	ss, _ := m.dbs.GetAllSlaves()
	for i, s := range ss {
		stats := s.Stats()
		opts := metric.WithAttributes(
			attribute.Key("db_name").String(fmt.Sprintf("slave_%d", i)),
		)

		o.ObserveFloat64(m.w, float64(stats.MaxLifetimeClosed), m.opts, opts)
	}

	return nil
}
