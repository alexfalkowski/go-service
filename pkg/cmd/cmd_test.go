package cmd_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/pkg/checker"
	"github.com/alexfalkowski/go-health/pkg/server"
	"github.com/alexfalkowski/go-service/pkg/cmd"
	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/health"
	healthGRPC "github.com/alexfalkowski/go-service/pkg/health/transport/grpc"
	healthHTTP "github.com/alexfalkowski/go-service/pkg/health/transport/http"
	"github.com/alexfalkowski/go-service/pkg/logger"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
)

func TestInvalidHTTP(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("APP_NAME", "test")
		os.Setenv("HTTP_PORT", "-1")
		os.Setenv("GRPC_PORT", "9000")
		os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I try to create a server", func() {
			opts := []fx.Option{
				logger.Module, pkgHTTP.Module, grpc.Module, config.Module,
			}

			err := cmd.RunServer([]string{}, 10*time.Second, opts)

			Convey("Then I should see an error", func() {
				So(err, ShouldBeError)
			})

			So(os.Unsetenv("APP_NAME"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("DATABASE_URL"), ShouldBeNil)
		})
	})
}

func TestInvalidGRPC(t *testing.T) {
	Convey("Given I have invalid HTTP port set", t, func() {
		os.Setenv("APP_NAME", "test")
		os.Setenv("HTTP_PORT", "9000")
		os.Setenv("GRPC_PORT", "-1")
		os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

		Convey("When I try to create a server", func() {
			opts := []fx.Option{
				logger.Module, pkgHTTP.Module, grpc.Module, config.Module,
				health.Module, fx.Provide(registrations), fx.Provide(httpObserver), fx.Provide(grpcObserver),
			}

			err := cmd.RunServer([]string{}, 10*time.Second, opts)

			Convey("Then I should see an error", func() {
				So(err, ShouldBeError)
			})

			So(os.Unsetenv("APP_NAME"), ShouldBeNil)
			So(os.Unsetenv("HTTP_PORT"), ShouldBeNil)
			So(os.Unsetenv("GRPC_PORT"), ShouldBeNil)
			So(os.Unsetenv("DATABASE_URL"), ShouldBeNil)
		})
	})
}

func registrations(cfg *config.Config, rtp http.RoundTripper) health.Registrations {
	hc := checker.NewHTTPCheckerWithRoundTripper("https://google.com", 1*time.Second, rtp)
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
