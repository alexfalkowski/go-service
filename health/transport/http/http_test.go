//nolint:varnamelen
package http_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-health/checker"
	"github.com/alexfalkowski/go-health/server"
	"github.com/alexfalkowski/go-service/v2/health"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
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

			test.RegisterHealth(o, o, o)
			world.RequireStart()

			Convey("When I query "+check, func() {
				ctx := t.Context()
				ctx = meta.WithRequestID(ctx, meta.String("test-id"))
				ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

				header := http.Header{}
				header.Set(content.TypeKey, mime.JSONMediaType)

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

		test.RegisterHealth(o, o, so.Observe("noop"))
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}
			header.Add("Request-Id", "test-id")
			header.Add("User-Agent", "test-user-agent")
			header.Set(content.TypeKey, mime.JSONMediaType)

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "readyz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a healthy response", func() {
				So(body, ShouldContainSubstring, "SERVING")
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get(content.TypeKey), ShouldEqual, mime.JSONMediaType)
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

		test.RegisterHealth(o, o, o)
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "healthz", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an unhealthy response", func() {
				So(body, ShouldEqual, "http: http checker: invalid status code")
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get(content.TypeKey), ShouldEqual, mime.ErrorMediaType)
			})

			world.RequireStop()
		})
	})
}

func observer(lc fx.Lifecycle, url string, world *test.World) (*server.Server, error) {
	dc, err := test.NewDBChecker(world)
	if err != nil {
		return nil, err
	}

	dr := server.NewRegistration("db", 10*time.Millisecond, dc)

	cc := checker.NewHTTPChecker(url, 5*time.Second, checker.WithRoundTripper(world.NewHTTP().Transport))
	hr := server.NewRegistration("http", 10*time.Millisecond, cc)

	no := checker.NewNoopChecker()
	nr := server.NewRegistration("noop", 10*time.Millisecond, no)

	regs := health.Registrations{hr, nr, dr}

	return health.NewServer(lc, regs), nil
}
