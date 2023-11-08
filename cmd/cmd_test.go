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
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/health"
	hgrpc "github.com/alexfalkowski/go-service/health/transport/grpc"
	hhttp "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/version"
	rcache "github.com/go-redis/cache/v8"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestShutdown(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			c := cmd.New()
			c.AddVersion("1.0.0")
			c.AddWorker(opts())

			Convey("Then I should not see an error", func() {
				So(c.RunWithArgs([]string{"worker"}), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestRun(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			c := cmd.New()
			c.AddVersion("1.0.0")
			c.AddWorker(opts())

			Convey("Then I should not see an error", func() {
				So(c.Run(), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalid(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		Convey("When I try to run an application", func() {
			c := cmd.New()
			c.AddServer(opts())

			Convey("Then I should see an error", func() {
				err := c.RunWithArgs([]string{"server", "--input", "file:../test/invalid.config.yml"})

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "invalid port")
			})
		})
	})
}

func TestClient(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		Convey("When I try to run a client", func() {
			opts := []fx.Option{fx.NopLogger}

			c := cmd.New()
			c.AddClient(opts)

			Convey("Then I should not see an error", func() {
				So(c.RunWithArgs([]string{"client"}), ShouldBeNil)
			})
		})
	})
}

func TestInvalidClient(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("TEST_CONFIG_FILE", "../test/invalid.config.yml")

		Convey("When I try to run an application", func() {
			c := cmd.New()
			c.AddClient(opts())

			Convey("Then I should see an error", func() {
				err := c.RunWithArgs([]string{"client", "--input", "env:TEST_CONFIG_FILE"})

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, "invalid port")
			})

			So(os.Unsetenv("TEST_CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func registrations(logger *zap.Logger, cfg *shttp.Config, tracer htracer.Tracer, _ version.Version) (health.Registrations, error) {
	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 5*time.Second, nc)

	client, err := shttp.NewClient(cfg,
		shttp.WithClientLogger(logger), shttp.WithClientTracer(tracer),
	)
	if err != nil {
		return nil, err
	}

	hc := checker.NewHTTPChecker("https://google.com", client)
	hr := server.NewRegistration("http", 5*time.Second, hc)

	return health.Registrations{nr, hr}, nil
}

func healthObserver(healthServer *server.Server) (*hhttp.HealthObserver, error) {
	return &hhttp.HealthObserver{Observer: healthServer.Observe("noop")}, nil
}

func livenessObserver(healthServer *server.Server) *hhttp.LivenessObserver {
	return &hhttp.LivenessObserver{Observer: healthServer.Observe("noop")}
}

func readinessObserver(healthServer *server.Server) *hhttp.ReadinessObserver {
	return &hhttp.ReadinessObserver{Observer: healthServer.Observe("http")}
}

func grpcObserver(healthServer *server.Server) *hgrpc.Observer {
	return &hgrpc.Observer{Observer: healthServer.Observe("http")}
}

func configs(c *rcache.Cache, _ *redis.Config, _ *ristretto.Config, _ *pg.Config, _ *nsq.Config) error {
	return c.Delete(context.Background(), "test")
}

func meter(_ metric.Meter) {
}

func ver() version.Version {
	return test.Version
}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(time.Second)

		_ = s.Shutdown()
	}(s)
}

func opts() []fx.Option {
	return []fx.Option{
		fx.NopLogger,
		runtime.Module, cmd.Module, config.Module, telemetry.Module, metrics.Module,
		health.Module, cache.RedisModule, cache.RistrettoModule,
		sql.PostgreSQLModule, transport.Module,
		cache.ProtoMarshallerModule, cache.SnappyCompressorModule,
		fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver),
		fx.Provide(readinessObserver), fx.Provide(grpcObserver), fx.Invoke(shutdown),
		fx.Invoke(configs), fx.Provide(ver), fx.Invoke(meter),
	}
}
