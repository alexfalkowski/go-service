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
	"github.com/linxGnu/mssqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// Driver is an alias for the driver.Driver type.
type Driver = driver.Driver

// Register registers a new driver.
func Register(name string, driver Driver) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	sql.Register(name, otelsql.WrapDriver(driver, otelsql.WithAttributes(semconv.DBSystemNameKey.String(name))))

	return err
}

// Open a DB pool.
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
	err := errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	attrs := otelsql.WithAttributes(semconv.DBSystemNameKey.String(name))

	masters, _ := db.GetAllMasters()
	for _, db := range masters {
		runtime.Must(otelsql.RegisterDBStatsMetrics(db.DB, attrs))
	}

	slaves, _ := db.GetAllSlaves()
	for _, db := range slaves {
		runtime.Must(otelsql.RegisterDBStatsMetrics(db.DB, attrs))
	}

	return db, nil
}
