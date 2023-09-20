package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/database/sql/config"
	sotel "github.com/alexfalkowski/go-service/database/sql/driver/otel"
	dzap "github.com/alexfalkowski/go-service/database/sql/driver/zap"
	"github.com/alexfalkowski/go-service/database/sql/metrics/prometheus"
	"github.com/alexfalkowski/go-service/database/sql/pg/otel"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/version"
	"github.com/linxGnu/mssqlx"
	"github.com/ngrok/sqlmw"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Register the driver for SQL.
func Register(name string, driver driver.Driver, tracer otel.Tracer, logger *zap.Logger) {
	var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}
	interceptor = sotel.NewInterceptor(name, tracer, interceptor)
	interceptor = dzap.NewInterceptor(name, logger, interceptor)

	sql.Register(name, sqlmw.Driver(driver, interceptor))
}

// Open a DB pool.
func Open(lc fx.Lifecycle, name string, cfg config.Config, ver version.Version) (*mssqlx.DBs, error) {
	masterDSNs := make([]string, len(cfg.Masters))
	for i, m := range cfg.Masters {
		masterDSNs[i] = m.URL
	}

	slaveDSNs := make([]string, len(cfg.Slaves))
	for i, s := range cfg.Slaves {
		slaveDSNs[i] = s.URL
	}

	db, err := connect(name, masterDSNs, slaveDSNs)
	if err != nil {
		return nil, err
	}

	prometheus.Register(lc, prometheus.NewCollector(name, db, ver))

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return destroy(db)
		},
	})

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func connect(name string, masterDSNs, slaveDSNs []string) (*mssqlx.DBs, error) {
	db, errs := mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)

	return db, errors.Combine(errs...)
}

func destroy(db *mssqlx.DBs) error {
	return errors.Combine(db.Destroy()...)
}
