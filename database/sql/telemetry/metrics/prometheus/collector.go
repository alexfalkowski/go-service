package prometheus

import (
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements the prometheus.Collector interface.
type Collector struct {
	db                *mssqlx.DBs
	maxOpen           *prometheus.Desc
	open              *prometheus.Desc
	inUse             *prometheus.Desc
	idle              *prometheus.Desc
	waitedFor         *prometheus.Desc
	blockedSeconds    *prometheus.Desc
	closedMaxIdle     *prometheus.Desc
	closedMaxLifetime *prometheus.Desc
	closedMaxIdleTime *prometheus.Desc
}

// NewCollector for prometheus.
func NewCollector(name string, db *mssqlx.DBs, version version.Version) *Collector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &Collector{
		db: db,
		maxOpen: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_max_open_total", name),
			fmt.Sprintf("Maximum number of open connections to %s.", name),
			[]string{"db_name"},
			labels,
		),
		open: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_open_total", name),
			fmt.Sprintf("The number of established connections both in use and idle for %s.", name),
			[]string{"db_name"},
			labels,
		),
		inUse: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_in_use_total", name),
			fmt.Sprintf("The number of connections currently in use for %s.", name),
			[]string{"db_name"},
			labels,
		),
		idle: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_idle_total", name),
			fmt.Sprintf("The number of idle connections for %s.", name),
			[]string{"db_name"},
			labels,
		),
		waitedFor: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_waited_for_total", name),
			fmt.Sprintf("The total number of connections waited for in %s.", name),
			[]string{"db_name"},
			labels,
		),
		blockedSeconds: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_blocked_seconds_total", name),
			fmt.Sprintf("The total time blocked waiting for a new connection for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxIdle: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_idle_total", name),
			fmt.Sprintf("The total number of connections closed due to SetMaxIdleConns for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxLifetime: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_lifetime_total", name),
			fmt.Sprintf("The total number of connections closed due to SetConnMaxLifetime for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxIdleTime: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_idle_time_total", name),
			fmt.Sprintf("The total number of connections closed due to SetConnMaxIdleTime for %s.", name),
			[]string{"db_name"},
			labels,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpen
	ch <- c.open
	ch <- c.inUse
	ch <- c.idle
	ch <- c.waitedFor
	ch <- c.blockedSeconds
	ch <- c.closedMaxIdle
	ch <- c.closedMaxLifetime
	ch <- c.closedMaxIdleTime
}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	ms, _ := c.db.GetAllMasters()
	for i, m := range ms {
		c.collect(fmt.Sprintf("master_%d", i), m, ch)
	}

	ss, _ := c.db.GetAllSlaves()
	for i, s := range ss {
		c.collect(fmt.Sprintf("slave_%d", i), s, ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (c Collector) collect(name string, db *sqlx.DB, ch chan<- prometheus.Metric) {
	stats := db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.maxOpen,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.open,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.inUse,
		prometheus.GaugeValue,
		float64(stats.InUse),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.idle,
		prometheus.GaugeValue,
		float64(stats.Idle),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.waitedFor,
		prometheus.CounterValue,
		float64(stats.WaitCount),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.blockedSeconds,
		prometheus.CounterValue,
		stats.WaitDuration.Seconds(),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxIdle,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxLifetime,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxIdleTime,
		prometheus.CounterValue,
		float64(stats.MaxIdleTimeClosed),
		name,
	)
}
