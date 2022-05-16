package driver

import (
	"database/sql"
	"database/sql/driver"

	"github.com/alexfalkowski/go-service/database/sql/driver/trace/opentracing"
	dzap "github.com/alexfalkowski/go-service/database/sql/driver/zap"
	"github.com/ngrok/sqlmw"
	otr "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// Register the driver for SQL.
func Register(name string, driver driver.Driver, tracer otr.Tracer, logger *zap.Logger) {
	var interceptor sqlmw.Interceptor = &sqlmw.NullInterceptor{}
	interceptor = opentracing.NewInterceptor(name, tracer, interceptor)
	interceptor = dzap.NewInterceptor(name, logger, interceptor)

	sql.Register(name, sqlmw.Driver(driver, interceptor))
}
