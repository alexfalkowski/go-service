package cmd_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/health"
	hgrpc "github.com/alexfalkowski/go-service/health/transport/grpc"
	hhttp "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/logger"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/sql"
	"github.com/alexfalkowski/go-service/sql/pg"
	"github.com/alexfalkowski/go-service/trace"
	"github.com/alexfalkowski/go-service/transport"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
	rcache "github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestShutdown(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, health.GRPCModule, health.HTTPModule, health.ServerModule,
				cache.RedisModule, cache.RistrettoModule, security.Auth0Module, sql.PostgreSQLModule,
				trace.DataDogOpenTracingModule, trace.JaegerOpenTracingModule,
				transport.HTTPServerModule, transport.GRPCServerModule,
				fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
				fx.Invoke(shutdown), fx.Invoke(configs),
			}

			c := cmd.New(10 * time.Second)
			c.AddWorker(opts)

			Convey("Then I should not see an error", func() {
				So(c.RunWithArg("worker"), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidHTTP(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../test/invalid_http.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
			}

			c := cmd.New(10 * time.Second)
			c.AddServer(opts)

			Convey("Then I should see an error", func() {
				err := c.RunWithArg("server")

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "listen tcp: address -1: invalid port")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalidGRPC(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("CONFIG_FILE", "../test/invalid_grpc.config.yml")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				fx.NopLogger,
				config.Module, logger.ZapModule, transport.HTTPServerModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(healthObserver), fx.Provide(livenessObserver), fx.Provide(readinessObserver), fx.Provide(grpcObserver),
			}

			c := cmd.New(10 * time.Second)
			c.AddServer(opts)

			Convey("Then I should see an error", func() {
				err := c.RunWithArg("server")

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "listen tcp: address -1: invalid port")
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestClient(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		Convey("When I try to run a client", func() {
			opts := []fx.Option{fx.NopLogger}

			c := cmd.New(10 * time.Second)
			c.AddClient(opts)

			Convey("Then I should not see an error", func() {
				So(c.RunWithArg("client"), ShouldBeNil)
			})
		})
	})
}

func registrations(logger *zap.Logger, cfg *shttp.Config) health.Registrations {
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 5*time.Second, nc)
	client := shttp.NewClient(cfg, logger)

	hc := checker.NewHTTPChecker("https://google.com", client)
	hr := server.NewRegistration("http", 5*time.Second, hc)

	return health.Registrations{nr, hr}
}

func healthObserver(healthServer *server.Server) (*hhttp.HealthObserver, error) {
	ob, err := healthServer.Observe("noop")
	if err != nil {
		return nil, err
	}

	return &hhttp.HealthObserver{Observer: ob}, nil
}

func livenessObserver(healthServer *server.Server) (*hhttp.LivenessObserver, error) {
	ob, err := healthServer.Observe("noop")
	if err != nil {
		return nil, err
	}

	return &hhttp.LivenessObserver{Observer: ob}, nil
}

func readinessObserver(healthServer *server.Server) (*hhttp.ReadinessObserver, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &hhttp.ReadinessObserver{Observer: ob}, nil
}

func grpcObserver(healthServer *server.Server) (*hgrpc.Observer, error) {
	ob, err := healthServer.Observe("http")
	if err != nil {
		return nil, err
	}

	return &hgrpc.Observer{Observer: ob}, nil
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
