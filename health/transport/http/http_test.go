//nolint:varnamelen
package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/health"
	shc "github.com/alexfalkowski/go-service/health/checker"
	shh "github.com/alexfalkowski/go-service/health/transport/http"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/meta"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
)

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()

			so, err := observer(world.Lifecycle, test.StatusURL("200"), world)
			So(err, ShouldBeNil)

			o := so.Observe("http")

			params := shh.RegisterParams{
				Health:   &shh.HealthObserver{Observer: o},
				Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
			}

			shh.Register(params)

			world.RequireStart()

			Convey("When I query "+check, func() {
				ctx := t.Context()
				ctx = tm.WithRequestID(ctx, meta.String("test-id"))
				ctx = tm.WithUserAgent(ctx, meta.String("test-user-agent"))

				header := http.Header{}
				header.Set("Content-Type", "application/json")

				res, body, err := world.ResponseWithBody(ctx, "http", world.InsecureServerHost(), http.MethodGet, check, header, http.NoBody)
				So(err, ShouldBeNil)

				Convey("Then I should have a healthy response", func() {
					So(res.StatusCode, ShouldEqual, http.StatusOK)
					So(body, ShouldContainSubstring, "SERVING")
				})

				world.RequireStop()
			})
		})
	}
}

func TestReadinessNoop(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()

		so, err := observer(world.Lifecycle, test.StatusURL("500"), world)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: so.Observe("noop")},
		}

		shh.Register(params)
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}
			header.Add("Request-Id", "test-id")
			header.Add("User-Agent", "test-user-agent")
			header.Set("Content-Type", "application/json")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "readyz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(body, ShouldContainSubstring, "SERVING")
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "application/json")
			})

			world.RequireStop()
		})
	})
}

func TestInvalidHealth(t *testing.T) {
	Convey("Given I register the health handler", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()

		so, err := observer(world.Lifecycle, test.StatusURL("500"), world)
		So(err, ShouldBeNil)

		o := so.Observe("http")

		params := shh.RegisterParams{
			Health:   &shh.HealthObserver{Observer: o},
			Liveness: &shh.LivenessObserver{Observer: o}, Readiness: &shh.ReadinessObserver{Observer: o},
		}

		shh.Register(params)
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "healthz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an unhealthy response", func() {
				So(body, ShouldEqual, "rest: http: invalid status code")
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.RequireStop()
		})
	})
}

func observer(lc fx.Lifecycle, url string, world *test.World) (*server.Server, error) {
	db, err := world.OpenDatabase()
	if err != nil {
		return nil, err
	}

	dc := shc.NewDBChecker(db, 1*time.Second)
	dr := server.NewRegistration("db", 10*time.Millisecond, dc)

	cc := checker.NewHTTPChecker(url, world.NewHTTP().Transport, 5*time.Second)
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, dr}

	return health.NewServer(lc, regs), nil
}
