package prometheus

import (
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	"github.com/prometheus/client_golang/prometheus"
)

// StatsCollector implements the prometheus.Collector interface.
type StatsCollector struct {
	db *mssqlx.DBs

	// descriptions of exported metrics
	maxOpenDesc           *prometheus.Desc
	openDesc              *prometheus.Desc
	inUseDesc             *prometheus.Desc
	idleDesc              *prometheus.Desc
	waitedForDesc         *prometheus.Desc
	blockedSecondsDesc    *prometheus.Desc
	closedMaxIdleDesc     *prometheus.Desc
	closedMaxLifetimeDesc *prometheus.Desc
	closedMaxIdleTimeDesc *prometheus.Desc
}

// NewStatsCollector for prometheus.
func NewStatsCollector(name string, db *mssqlx.DBs, version version.Version) *StatsCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	return &StatsCollector{
		db: db,
		maxOpenDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_max_open_total", name),
			fmt.Sprintf("Maximum number of open connections to %s.", name),
			[]string{"db_name"},
			labels,
		),
		openDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_open_total", name),
			fmt.Sprintf("The number of established connections both in use and idle for %s.", name),
			[]string{"db_name"},
			labels,
		),
		inUseDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_in_use_total", name),
			fmt.Sprintf("The number of connections currently in use for %s.", name),
			[]string{"db_name"},
			labels,
		),
		idleDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_idle_total", name),
			fmt.Sprintf("The number of idle connections for %s.", name),
			[]string{"db_name"},
			labels,
		),
		waitedForDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_waited_for_total", name),
			fmt.Sprintf("The total number of connections waited for in %s.", name),
			[]string{"db_name"},
			labels,
		),
		blockedSecondsDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_blocked_seconds_total", name),
			fmt.Sprintf("The total time blocked waiting for a new connection for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxIdleDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_idle_total", name),
			fmt.Sprintf("The total number of connections closed due to SetMaxIdleConns for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxLifetimeDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_lifetime_total", name),
			fmt.Sprintf("The total number of connections closed due to SetConnMaxLifetime for %s.", name),
			[]string{"db_name"},
			labels,
		),
		closedMaxIdleTimeDesc: prometheus.NewDesc(
			fmt.Sprintf("%s_sql_closed_max_idle_time_total", name),
			fmt.Sprintf("The total number of connections closed due to SetConnMaxIdleTime for %s.", name),
			[]string{"db_name"},
			labels,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenDesc
	ch <- c.openDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitedForDesc
	ch <- c.blockedSecondsDesc
	ch <- c.closedMaxIdleDesc
	ch <- c.closedMaxLifetimeDesc
	ch <- c.closedMaxIdleTimeDesc
}

// Collect implements the prometheus.Collector interface.
func (c StatsCollector) Collect(ch chan<- prometheus.Metric) {
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
func (c StatsCollector) collect(name string, db *sqlx.DB, ch chan<- prometheus.Metric) {
	stats := db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.maxOpenDesc,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.openDesc,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.inUseDesc,
		prometheus.GaugeValue,
		float64(stats.InUse),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.idleDesc,
		prometheus.GaugeValue,
		float64(stats.Idle),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.waitedForDesc,
		prometheus.CounterValue,
		float64(stats.WaitCount),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.blockedSecondsDesc,
		prometheus.CounterValue,
		stats.WaitDuration.Seconds(),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxIdleDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxLifetimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
		name,
	)
	ch <- prometheus.MustNewConstMetric(
		c.closedMaxIdleTimeDesc,
		prometheus.CounterValue,
		float64(stats.MaxIdleTimeClosed),
		name,
	)
}
