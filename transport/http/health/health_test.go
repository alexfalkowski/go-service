package health_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()

			server := world.HealthServer(test.Name.String(), test.StatusURL("200"))

			err := server.Observe(test.Name.String(), check, "http")
			So(err, ShouldBeNil)

			test.RegisterHealth(server)
			world.RequireStart()

			Convey("When I query "+check, func() {
				ctx := t.Context()
				ctx = meta.WithRequestID(ctx, meta.String("test-id"))
				ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

				header := http.Header{}
				header.Set(content.TypeKey, mime.JSONMediaType)

				url := world.NamedServerURL("http", check)

				res, body, err := world.ResponseWithBody(ctx, url, http.MethodGet, header, http.NoBody)
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

		server := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := server.Observe(test.Name.String(), "readyz", "noop")
		So(err, ShouldBeNil)

		test.RegisterHealth(server)
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}
			header.Add("Request-Id", "test-id")
			header.Add("User-Agent", "test-user-agent")
			header.Set(content.TypeKey, mime.JSONMediaType)

			url := world.NamedServerURL("http", "readyz")

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
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

		server := world.HealthServer(test.Name.String(), test.StatusURL("500"))

		err := server.Observe(test.Name.String(), "healthz", "http")
		So(err, ShouldBeNil)

		test.RegisterHealth(server)
		world.RequireStart()

		Convey("When I query health", func() {
			header := http.Header{}
			url := world.NamedServerURL("http", "healthz")

			res, body, err := world.ResponseWithBody(t.Context(), url, http.MethodGet, header, http.NoBody)
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

func TestMissingHealth(t *testing.T) {
	checks := []string{"healthz", "livez", "readyz"}

	for _, check := range checks {
		Convey("Given I register the health handler", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()

			server := world.HealthServer(test.Name.String(), test.StatusURL("200"))

			test.RegisterHealth(server)
			world.RequireStart()

			Convey("When I query "+check, func() {
				ctx := t.Context()
				ctx = meta.WithRequestID(ctx, meta.String("test-id"))
				ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

				header := http.Header{}
				header.Set(content.TypeKey, mime.JSONMediaType)

				url := world.NamedServerURL("http", check)

				res, err := world.ResponseWithNoBody(ctx, url, http.MethodGet, header)
				So(err, ShouldBeNil)

				Convey("Then I should have a unhealthy response", func() {
					So(res.StatusCode, ShouldEqual, http.StatusServiceUnavailable)
				})

				world.RequireStop()
			})
		})
	}
}
