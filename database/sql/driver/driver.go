package driver

import (
	"database/sql"
	"database/sql/driver"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/telemetry"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
)

// Driver aliases [database/sql/driver.Driver].
//
// It is the concrete driver type expected by Register.
type Driver = driver.Driver

// Conn aliases [database/sql/driver.Conn].
type Conn = driver.Conn

// NamedValue aliases [database/sql/driver.NamedValue].
type NamedValue = driver.NamedValue

// Rows aliases [database/sql/driver.Rows].
type Rows = driver.Rows

// Stmt aliases [database/sql/driver.Stmt].
type Stmt = driver.Stmt

// Tx aliases [database/sql/driver.Tx].
type Tx = driver.Tx

// Value aliases [database/sql/driver.Value].
type Value = driver.Value

// ErrSkip aliases [database/sql/driver.ErrSkip].
var ErrSkip = driver.ErrSkip

// ErrNoDSNs is returned when SQL configuration enables a driver without any master or slave DSNs.
var ErrNoDSNs = errors.New("driver: no database DSNs configured")

// ErrEmptyDSN is returned when a configured DSN source resolves to an empty string.
var ErrEmptyDSN = errors.New("driver: empty database DSN")

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
		driver = telemetry.WrapDriver(driver, telemetryOptions(name, opts)...)
	}

	sql.Register(name, driver)

	return err
}

// Connect opens master/slave [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a previously
// registered driver name.
//
// It resolves DSNs from cfg using the provided filesystem (DSNs are configured
// as go-service "source strings"), connects using
// [mssqlx.ConnectMasterSlaves], registers OpenTelemetry DB stats metrics for
// each pool when metrics are enabled, and then applies pool settings
// (connection lifetime, max idle, and max open connections).
//
// Preconditions:
//   - cfg must be non-nil and already treated as enabled/validated by the caller.
//
// Failure behavior:
//   - returns errors encountered while resolving DSNs or connecting,
//   - returns [ErrEmptyDSN] when any configured DSN source resolves to empty bytes, and
//   - returns [ErrNoDSNs] when neither masters nor slaves are configured.
//
// The returned DBs embeds the upstream master/slave pool collection and owns
// repository lifecycle cleanup such as DB stats metric unregistration.
func Connect(name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	return connect(name, fs, cfg, opts...)
}

// Open opens master/slave [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a previously
// registered driver name.
//
// Open delegates the connection work to [Connect] and then appends an OnStop hook
// to the provided lifecycle that closes all returned pools by calling [DBs.Destroy].
//
// Preconditions:
//   - cfg must be non-nil and already treated as enabled/validated by the caller.
//   - driver-specific wrappers, such as [github.com/alexfalkowski/go-service/v2/database/sql/pg.Open], own nil/disabled
//     config semantics before delegating here.
//
// The returned type is the same go-service [DBs] wrapper returned by [Connect].
func Open(lc di.Lifecycle, name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	db, err := Connect(name, fs, cfg, opts...)
	if err != nil {
		return nil, err
	}

	lc.Append(di.Hook{
		OnStop: func(_ context.Context) error {
			return db.Destroy()
		},
	})

	return db, nil
}

// ConnectMasterSlaves opens master/slave [github.com/alexfalkowski/go-service/v2/database/sql] connection pools for a
// previously registered driver name.
//
// It is a low-level wrapper around [mssqlx.ConnectMasterSlaves] for callers that
// already have resolved literal DSNs.
func ConnectMasterSlaves(name string, masterDSNs, slaveDSNs []string) (*DBs, []error) {
	db, errs := mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)
	if errors.Join(errs...) != nil {
		if db != nil {
			errs = append(errs, db.Destroy()...)
		}

		return nil, errs
	}

	return &DBs{DBs: db}, nil
}

func connect(name string, fs *os.FS, cfg *config.Config, opts ...telemetry.Option) (*DBs, error) {
	masters, err := resolveDSNs(fs, cfg.Masters)
	if err != nil {
		return nil, err
	}

	slaves, err := resolveDSNs(fs, cfg.Slaves)
	if err != nil {
		return nil, err
	}

	db, err := connectDBs(name, masters, slaves, opts...)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime.Duration())
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func resolveDSNs(fs *os.FS, dsns []config.DSN) ([]string, error) {
	resolved := make([]string, len(dsns))

	for i, dsn := range dsns {
		url, err := dsn.GetURL(fs)
		if err != nil {
			return nil, err
		}
		if len(url) == 0 {
			return nil, ErrEmptyDSN
		}

		resolved[i] = bytes.String(url)
	}

	return resolved, nil
}

func connectDBs(name string, masterDSNs, slaveDSNs []string, opts ...telemetry.Option) (*DBs, error) {
	if len(masterDSNs)+len(slaveDSNs) == 0 {
		return nil, ErrNoDSNs
	}

	db, errs := mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)
	if err := errors.Join(errs...); err != nil {
		if db != nil {
			err = errors.Join(err, errors.Join(db.Destroy()...))
		}

		return nil, err
	}

	var regs []metrics.Registration
	if metrics.IsEnabled() {
		opts := telemetryOptions(name, opts)

		masters, _ := db.GetAllMasters()
		regs = append(regs, register(masters, "master", opts...)...)

		slaves, _ := db.GetAllSlaves()
		regs = append(regs, register(slaves, "slave", opts...)...)
	}

	return &DBs{DBs: db, registrations: regs}, nil
}

func register(dbs []*sqlx.DB, role string, opts ...telemetry.Option) []metrics.Registration {
	regs := make([]metrics.Registration, 0, len(dbs))

	for i, db := range dbs {
		reg, err := telemetry.RegisterDBStatsMetrics(db.DB, dbStatsOptions(role, i, opts)...)
		runtime.Must(err)
		regs = append(regs, reg)
	}

	return regs
}

func dbStatsOptions(role string, index int, opts []telemetry.Option) []telemetry.Option {
	options := make([]telemetry.Option, 0, len(opts)+1)
	options = append(options, opts...)

	// mssqlx does not expose a pool name, so DB stats metrics use a stable
	// go-service pattern unique within this DBs collection.
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

func telemetryOptions(name string, opts []telemetry.Option) []telemetry.Option {
	if len(opts) > 0 {
		return opts
	}

	return []telemetry.Option{telemetry.WithAttributes(attributes.DBSystemNameKey.String(name))}
}
