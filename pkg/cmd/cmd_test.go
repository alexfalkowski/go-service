package cmd_test

import (
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/cache"
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/cmd"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthGRPC "github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	healthHTTP "github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"github.com/alexfalkowski/go-service/pkg/logger"
	"github.com/alexfalkowski/go-service/pkg/security"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/sql"
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	"github.com/alexfalkowski/go-service/pkg/trace"
	"github.com/alexfalkowski/go-service/pkg/transport"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestShutdown(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			opts := []fx.Option{
				config.Module, logger.ZapModule, health.GRPCModule, health.HTTPModule, health.ServerModule,
				cache.RedisModule, cache.RistrettoModule, security.Auth0Module, sql.PostgreSQLModule,
				trace.DataDogOpenTracingModule, trace.JaegerOpenTracingModule,
				transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule, transport.NSQModule,
				fx.Provide(registrations), fx.Provide(httpObserver), fx.Provide(grpcObserver), fx.Invoke(shutdown),
				fx.Invoke(configs),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"worker"})

			Convey("Then I should not see an error", func() {
				So(c.Execute(), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidHTTP(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/invalid_http.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(httpObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				So(c.Execute(), ShouldBeError)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidGRPC(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/invalid_grpc.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(httpObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				So(c.Execute(), ShouldBeError)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func registrations(logger *zap.Logger) health.Registrations {
	hc := checker.NewHTTPCheckerWithClient("https://google.com", 1*time.Second, pkgHTTP.NewClient(logger))
	hr := server.NewRegistration("http", 5*time.Second, hc)

	return health.Registrations{hr}
}

func httpObserver(healthServer *server.Server) (*healthHTTP.Observer, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &healthHTTP.Observer{Observer: ob}, nil
}

func grpcObserver(healthServer *server.Server) (*healthGRPC.Observer, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &healthGRPC.Observer{Observer: ob}, nil
}

func configs(_ *redis.Config, _ *ristretto.Config, _ *auth0.Config, _ *pg.Config, _ *nsq.Config) {
}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(5 * time.Second)

		s.Shutdown() // nolint:errcheck
	}(s)
}
