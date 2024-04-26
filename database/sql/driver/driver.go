package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/database/sql/config"
	logger "github.com/alexfalkowski/go-service/database/sql/driver/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/database/sql/driver/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/linxGnu/mssqlx"
	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// Register the driver for SQL.
func Register(name string, driver driver.Driver, t trace.Tracer, l *zap.Logger) {
	var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}

	if t != nil {
		interceptor = tracer.NewInterceptor(name, t, interceptor)
	}

	if l != nil {
		interceptor = logger.NewInterceptor(name, l, interceptor)
	}

	sql.Register(name, sqlmw.Driver(driver, interceptor))
}

// Open a DB pool.
func Open(lc fx.Lifecycle, name string, cfg config.Config) (*mssqlx.DBs, error) {
	masters := make([]string, len(cfg.Masters))
	for i, m := range cfg.Masters {
		masters[i] = m.URL
	}

	slaves := make([]string, len(cfg.Slaves))
	for i, s := range cfg.Slaves {
		slaves[i] = s.URL
	}

	db, err := connect(name, masters, slaves)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return destroy(db)
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

	return db, multierr.Combine(errs...)
}

func destroy(db *mssqlx.DBs) error {
	return multierr.Combine(db.Destroy()...)
}
