package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/alexfalkowski/go-service/database/sql/config"
	dzap "github.com/alexfalkowski/go-service/database/sql/driver/telemetry/logger/zap"
	stracer "github.com/alexfalkowski/go-service/database/sql/driver/telemetry/tracer"
	"github.com/linxGnu/mssqlx"
	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// Register the driver for SQL.
func Register(name string, driver driver.Driver, tracer trace.Tracer, logger *zap.Logger) {
	var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}

	if tracer != nil {
		interceptor = stracer.NewInterceptor(name, tracer, interceptor)
	}

	if logger != nil {
		interceptor = dzap.NewInterceptor(name, logger, interceptor)
	}

	sql.Register(name, sqlmw.Driver(driver, interceptor))
}

// Open a DB pool.
func Open(lc fx.Lifecycle, name string, cfg config.Config) (*mssqlx.DBs, error) {
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

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return destroy(db)
		},
	})

	d, err := time.ParseDuration(cfg.ConnMaxLifetime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(d)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func connect(name string, masterDSNs, slaveDSNs []string) (*mssqlx.DBs, error) {
	db, errs := mssqlx.ConnectMasterSlaves(name, masterDSNs, slaveDSNs)

	return db, multierr.Combine(errs...)
}

func destroy(db *mssqlx.DBs) error {
	return multierr.Combine(db.Destroy()...)
}
