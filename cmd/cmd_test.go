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
	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/config"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/health"
	shg "github.com/alexfalkowski/go-service/health/transport/grpc"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/marshaller"
	sr "github.com/alexfalkowski/go-service/ristretto"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/time/ntp"
	"github.com/alexfalkowski/go-service/time/nts"
	"github.com/alexfalkowski/go-service/transport"
	geh "github.com/alexfalkowski/go-service/transport/events/http"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/version"
	v8 "github.com/go-redis/cache/v8"
	"github.com/open-feature/go-sdk/openfeature"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestRunWithServer(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/configs/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			c := cmd.New("1.0.0")
			c.AddServer(opts()...)

			Convey("Then I should not see an error", func() {
				So(c.RunWithArgs([]string{"server"}), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestRun(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("CONFIG_FILE", "../test/configs/config.yml")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			c := cmd.New("1.0.0")
			c.AddServer(opts()...)

			Convey("Then I should not see an error", func() {
				So(c.Run(), ShouldBeNil)
			})

			So(os.Unsetenv("CONFIG_FILE"), ShouldBeNil)
		})
	})
}

func TestInvalid(t *testing.T) {
	configs := []string{
		"file:../test/configs/invalid_http.config.yml",
		"file:../test/configs/invalid_grpc.config.yml",
		"file:../test/configs/invalid_debug.config.yml",
	}

	for _, i := range configs {
		Convey("Given I have an invalid configuration", t, func() {
			Convey("When I try to run an application", func() {
				c := cmd.New("1.0.0")
				c.AddServer(opts()...)

				Convey("Then I should not see an error", func() {
					err := c.RunWithArgs([]string{"server", "--input", i})

					So(err, ShouldBeError)
					So(err.Error(), ShouldContainSubstring, "invalid port")
				})
			})
		})
	}
}

func TestDisabled(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		Convey("When I try to run an application", func() {
			c := cmd.New("1.0.0")
			c.AddServer(opts()...)

			Convey("Then I should see an error", func() {
				err := c.RunWithArgs([]string{"server", "-i", "file:../test/configs/disabled.config.yml"})

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestClient(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		Convey("When I try to run a client", func() {
			opts := []fx.Option{fx.NopLogger}

			c := cmd.New("1.0.0")
			c.AddClient(opts...)

			Convey("Then I should not see an error", func() {
				So(c.RunWithArgs([]string{"client"}), ShouldBeNil)
			})
		})
	})
}

func TestInvalidClient(t *testing.T) {
	configs := []string{
		"../test/configs/invalid_http.config.yml",
		"../test/configs/invalid_grpc.config.yml",
	}

	for _, i := range configs {
		Convey("Given I have invalid configuration", t, func() {
			os.Setenv("TEST_CONFIG_FILE", i)

			Convey("When I try to run an application", func() {
				c := cmd.New("1.0.0")
				c.AddClient(opts()...)

				Convey("Then I should see an error", func() {
					err := c.RunWithArgs([]string{"client", "--input", "env:TEST_CONFIG_FILE"})

					So(err, ShouldBeError)
					So(err.Error(), ShouldContainSubstring, "invalid port")
				})

				So(os.Unsetenv("TEST_CONFIG_FILE"), ShouldBeNil)
			})
		})
	}
}

func registrations(logger *zap.Logger, cfg *http.Config, tracer trace.Tracer, _ version.Version) health.Registrations {
	if cfg == nil {
		return nil
	}

	nc := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 5*time.Second, nc)

	client := http.NewClient(http.WithClientLogger(logger), http.WithClientTracer(tracer), http.WithClientUserAgent(cfg.UserAgent))
	hc := checker.NewHTTPChecker("https://google.com", client)
	hr := server.NewRegistration("http", 5*time.Second, hc)

	return health.Registrations{nr, hr}
}

func healthObserver(healthServer *server.Server) (*shh.HealthObserver, error) {
	return &shh.HealthObserver{Observer: healthServer.Observe("noop")}, nil
}

func livenessObserver(healthServer *server.Server) *shh.LivenessObserver {
	return &shh.LivenessObserver{Observer: healthServer.Observe("noop")}
}

func readinessObserver(healthServer *server.Server) *shh.ReadinessObserver {
	return &shh.ReadinessObserver{Observer: healthServer.Observe("http")}
}

func grpcObserver(healthServer *server.Server) *shg.Observer {
	return &shg.Observer{Observer: healthServer.Observe("http")}
}

func redisCache(c *v8.Cache) error {
	return c.Delete(context.Background(), "test")
}

func ristrettoCache(c sr.Cache) {
	c.Del("test")
}

func configs(_ *redis.Config, _ *ristretto.Config, _ *pg.Config, _ *token.Config) {
}

func meter(_ metric.Meter) {
}

func featureClient(_ *openfeature.Client) {
}

func webHooks(_ *h.Webhook, _ *geh.Receiver) {
}

func ver() version.Version {
	return test.Version
}

func timeServices(nt *ntp.Service, ns *nts.Service) {
	nt.Query()
	nt.Time()

	ns.Query()
}

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(time.Second)

		_ = s.Shutdown()
	}(s)
}

func opts() []fx.Option {
	tm := fx.Options(
		transport.Module,
		fx.Provide(grpc.UnaryServerInterceptor),
		fx.Provide(grpc.StreamServerInterceptor),
		fx.Provide(http.ServerHandlers),
	)

	return []fx.Option{
		fx.NopLogger,
		runtime.Module, cmd.Module, config.Module, debug.Module, feature.Module, st.Module,
		telemetry.Module, metrics.Module, health.Module, sql.Module, tm, hooks.Module,
		cache.Module, compressor.Module, marshaller.Module,
		fx.Provide(registrations), fx.Provide(healthObserver), fx.Provide(livenessObserver),
		fx.Provide(readinessObserver), fx.Provide(grpcObserver), fx.Invoke(shutdown),
		fx.Invoke(featureClient), fx.Invoke(webHooks), fx.Invoke(configs),
		fx.Invoke(redisCache), fx.Invoke(ristrettoCache),
		fx.Provide(ver), fx.Invoke(meter), fx.Invoke(timeServices),
	}
}
