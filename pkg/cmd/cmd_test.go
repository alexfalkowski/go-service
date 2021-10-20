package cmd_test

import (
	"context"
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
	rcache "github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestShutdown(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, health.GRPCModule, health.HTTPModule, health.ServerModule,
				cache.RedisModule, cache.RistrettoModule, security.Auth0Module, sql.PostgreSQLModule,
				trace.DataDogOpenTracingModule, trace.JaegerOpenTracingModule,
				transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule, transport.NSQModule,
				fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
				fx.Invoke(shutdown), fx.Invoke(configs),
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

// nolint:dupl
func TestInvalidHTTP(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/invalid_http.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				err := c.Execute()

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "listen tcp: address -1: invalid port")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

// nolint:dupl
func TestInvalidGRPC(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../../test/invalid_grpc.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				err := c.Execute()

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "listen tcp: address -1: invalid port")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func registrations(logger *zap.Logger) health.Registrations {
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 5*time.Second, nc)

	hc := checker.NewHTTPCheckerWithClient("https://google.com", 1*time.Second, pkgHTTP.NewClient(logger))
	hr := server.NewRegistration("http", 5*time.Second, hc)

	return health.Registrations{nr, hr}
}

func healthObserver(healthServer *server.Server) (*healthHTTP.HealthObserver, error) {
	ob, err := healthServer.Observe("noop")
	if err != nil {
		return nil, err
	}

	return &healthHTTP.HealthObserver{Observer: ob}, nil
}

func livenessObserver(healthServer *server.Server) (*healthHTTP.LivenessObserver, error) {
	ob, err := healthServer.Observe("noop")
	if err != nil {
		return nil, err
	}

	return &healthHTTP.LivenessObserver{Observer: ob}, nil
}

func readinessObserver(healthServer *server.Server) (*healthHTTP.ReadinessObserver, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &healthHTTP.ReadinessObserver{Observer: ob}, nil
}

func grpcObserver(healthServer *server.Server) (*healthGRPC.Observer, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &healthGRPC.Observer{Observer: ob}, nil
}

func configs(c *rcache.Cache, _ *redis.Config, _ *ristretto.Config, _ *auth0.Config, _ *pg.Config, _ *nsq.Config) error {
	return c.Delete(context.Background(), "test")
}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(5 * time.Second)

		s.Shutdown() // nolint:errcheck
	}(s)
}
