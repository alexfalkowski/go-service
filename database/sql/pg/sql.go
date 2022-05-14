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

// DBParams for SQL.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Tracer    opentracing.Tracer
	Logger    *zap.Logger
	Version   version.Version
}

var once sync.Once

// NewDB for SQL.
func NewDB(params DBParams) *sql.DB {
	once.Do(func() {
		var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}
		interceptor = opentracing.NewInterceptor(params.Tracer, interceptor)
		interceptor = szap.NewInterceptor(params.Logger, interceptor)

		sql.Register(driverName, sqlmw.Driver(stdlib.GetDefaultDriver(), interceptor))
	})

	db, _ := sql.Open(driverName, params.Config.URL)

	prometheus.Register(params.Lifecycle, db, params.Version)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
