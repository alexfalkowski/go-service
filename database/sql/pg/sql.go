package pg

import (
	"context"
	"database/sql"

	"github.com/alexfalkowski/go-service/database/sql/metrics/prometheus"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/ngrok/sqlmw"
	"go.uber.org/fx"
)

const driverName = "pgx-mw"

// DBParams for SQL.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Tracer    opentracing.Tracer
	Version   version.Version
}

// NewDB for SQL.
func NewDB(params DBParams) *sql.DB {
	if !isDriverRegistered() {
		var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}
		interceptor = opentracing.NewInterceptor(params.Tracer, interceptor)

		sql.Register(driverName, sqlmw.Driver(stdlib.GetDefaultDriver(), interceptor))
	}

	db, _ := sql.Open(driverName, params.Config.URL)

	prometheus.Register(params.Lifecycle, db, params.Version)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}

// isDriverRegistered as we can only really register it once. This is only a problem during tests.
func isDriverRegistered() bool {
	for _, d := range sql.Drivers() {
		if d == driverName {
			return true
		}
	}

	return false
}
