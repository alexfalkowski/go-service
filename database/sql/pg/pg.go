package pg

import (
	"context"
	"database/sql"
	"sync"

	"github.com/alexfalkowski/go-service/database/sql/metrics/prometheus"
	szap "github.com/alexfalkowski/go-service/database/sql/pg/logger/zap"
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/ngrok/sqlmw"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const driverName = "pgx-mw"

// DBParams for PostgreSQL.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Version   version.Version
}

// Open for PostgreSQL.
func Open(params DBParams) *sql.DB {
	db, _ := sql.Open(driverName, params.Config.URL)

	prometheus.Register(params.Lifecycle, db, params.Version)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}

var once sync.Once

// Register the driver for PostgreSQL.
func Register(tracer opentracing.Tracer, logger *zap.Logger) {
	once.Do(func() {
		var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}
		interceptor = opentracing.NewInterceptor(tracer, interceptor)
		interceptor = szap.NewInterceptor(logger, interceptor)

		sql.Register(driverName, sqlmw.Driver(stdlib.GetDefaultDriver(), interceptor))
	})
}
