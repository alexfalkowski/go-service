package cmd_test

import (
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/cmd"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthGRPC "github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	healthHTTP "github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"github.com/alexfalkowski/go-service/pkg/logger"
	"github.com/alexfalkowski/go-service/pkg/transport"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestShutdown(t *testing.T) {
	Convey("Given I have valid configuration", t, func() {
		os.Setenv("SERVICE_NAME", "test")
		os.Setenv("SERVICE_DESCRIPTION", "Test service.")
		os.Setenv("HTTP_PORT", "8000")
		os.Setenv("GRPC_PORT", "9000")
		os.Setenv("POSTGRESQL_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I try to run an application that will shutdown in 5 seconds", func() {
			opts := []fx.Option{
				logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(httpObserver), fx.Provide(grpcObserver), fx.Invoke(shutdown),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"worker"})

			Convey("Then I should not see an error", func() {
				So(c.Execute(), ShouldBeNil)
			})

			So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
			So(os.Unsetenv("SERVICE_DESCRIPTION"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("POSTGRESQL_URL"), ShouldBeNil)
		})
	})
}

// nolint:dupl
func TestInvalidHTTP(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("SERVICE_NAME", "test")
		os.Setenv("SERVICE_DESCRIPTION", "Test service.")
		os.Setenv("HTTP_PORT", "-1")
		os.Setenv("GRPC_PORT", "9000")
		os.Setenv("POSTGRESQL_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(httpObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				So(c.Execute(), ShouldBeError)
			})

			So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
			So(os.Unsetenv("SERVICE_DESCRIPTION"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("POSTGRESQL_URL"), ShouldBeNil)
		})
	})
}

// nolint:dupl
func TestInvalidGRPC(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("SERVICE_NAME", "test")
		os.Setenv("SERVICE_DESCRIPTION", "Test service.")
		os.Setenv("HTTP_PORT", "9000")
		os.Setenv("GRPC_PORT", "-1")
		os.Setenv("POSTGRESQL_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I try to run an application", func() {
			opts := []fx.Option{
				logger.ZapModule, transport.HTTPServerModule, transport.HTTPClientModule, transport.GRPCServerModule,
				health.GRPCModule, health.HTTPModule, health.ServerModule, fx.Provide(registrations),
				fx.Provide(httpObserver), fx.Provide(grpcObserver),
			}

			c, err := cmd.New(10*time.Second, opts, opts)
			So(err, ShouldBeNil)

			c.SetArgs([]string{"serve"})

			Convey("Then I should see an error", func() {
				So(c.Execute(), ShouldBeError)
			})

			So(os.Unsetenv("SERVICE_NAME"), ShouldBeNil)
			So(os.Unsetenv("SERVICE_DESCRIPTION"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("POSTGRESQL_URL"), ShouldBeNil)
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

func shutdown(s fx.Shutdowner) {
	go func(s fx.Shutdowner) {
		time.Sleep(5 * time.Second)

		s.Shutdown() // nolint:errcheck
	}(s)
}
