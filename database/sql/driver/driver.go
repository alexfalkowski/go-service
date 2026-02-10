package driver

import (
	"database/sql"
	"database/sql/driver"

	"github.com/XSAM/otelsql"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/jmoiron/sqlx"
	"github.com/linxGnu/mssqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// Driver aliases `database/sql/driver`.Driver.
type Driver = driver.Driver

// Register registers a `database/sql` driver under name and wraps it with OpenTelemetry
// instrumentation.
//
// If the underlying `sql.Register` panics (for example, due to a duplicate name),
// Register converts the panic into an error and returns it.
func Register(name string, driver Driver) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	sql.Register(name, otelsql.WrapDriver(driver, otelsql.WithAttributes(semconv.DBSystemNameKey.String(name))))

	return err
}

// Open opens master/slave `database/sql` connection pools for the given driver name.
//
// It reads master and slave DSNs from cfg using the provided filesystem (DSNs are
// configured as "source strings"), connects using `mssqlx.ConnectMasterSlaves`,
// registers OpenTelemetry DB stats metrics for each pool, and configures pool
// limits/lifetime.
//
// The returned pools are closed on Fx lifecycle stop.
func Open(lc di.Lifecycle, name string, fs *os.FS, cfg *config.Config) (*mssqlx.DBs, error) {
	masters := make([]string, len(cfg.Masters))

	for i, m := range cfg.Masters {
		u, err := m.GetURL(fs)
		if err != nil {
			return nil, err
		}

		masters[i] = bytes.String(u)
	}

	slaves := make([]string, len(cfg.Slaves))

	for i, s := range cfg.Slaves {
		u, err := s.GetURL(fs)
		if err != nil {
			return nil, err
		}

		slaves[i] = bytes.String(u)
	}

	db, err := connect(name, masters, slaves)
	if err != nil {
		return nil, err
	}

	lc.Append(di.Hook{
		OnStop: func(_ context.Context) error {
			return errors.Join(db.Destroy()...)
		},
	})

	d := time.MustParseDuration(cfg.ConnMaxLifetime)

	db.SetConnMaxLifetime(d)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func connect(name string, masterDSNs, slaveDSNs []string) (*mssqlx.DBs, error) {
	db, errs := mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	attrs := otelsql.WithAttributes(semconv.DBSystemNameKey.String(name))

	masters, _ := db.GetAllMasters()
	register(masters, attrs)

	slaves, _ := db.GetAllSlaves()
	register(slaves, attrs)

	return db, nil
}

func register(dbs []*sqlx.DB, opts ...otelsql.Option) {
	for _, db := range dbs {
		_, err := otelsql.RegisterDBStatsMetrics(db.DB, opts...)
		runtime.Must(err)
	}
}
