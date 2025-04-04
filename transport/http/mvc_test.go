package http_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/status"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/html"
)

func TestRouteSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
			return mvc.View("hello.tmpl"), &test.Model, nil
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	})
}

func TestRouteMissingView(t *testing.T) {
	Convey("Given I have setup a route with an missisng view", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
			return mvc.View("none.tmpl"), &test.Model, nil
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have no html", func() {
				So(body, ShouldBeEmpty)
				So(res.StatusCode, ShouldEqual, 500)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")
			})

			world.RequireStop()
		})
	})
}

func TestRouteError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
			return mvc.View("error.tmpl"), &test.Model, status.FromError(http.StatusServiceUnavailable, test.ErrFailed)
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	})
}

func TestStaticSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.Static("GET /robots.txt", "static/robots.txt")

		Convey("When I query for robots", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.RequireStop()
		})
	})
}

func TestStaticError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.Static("GET /robots.txt", "static/bob.txt")

		Convey("When I query for hello", func() {
			header := http.Header{}

			res, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(res.StatusCode, ShouldEqual, 500)
			})

			world.RequireStop()
		})
	})
}

func TestMissingViews(t *testing.T) {
	Convey("Given I have views with missing FS", t, func() {
		mux := http.NewServeMux()

		v := mvc.NewViews(mvc.ViewsParams{Patterns: mvc.Patterns{"views/*.tmpl"}})
		mvc.Register(mux, v)

		Convey("When I add routes to an invalid view", func() {
			routeAdded := mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
				return mvc.View("hello.tmpl"), &test.Model, nil
			})
			staticAdded := mvc.Static("GET /robots.txt", "static/bob.txt")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(staticAdded, ShouldBeFalse)
			})
		})
	})

	Convey("Given I have views with missing patterns", t, func() {
		mux := http.NewServeMux()
		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views})

		mvc.Register(mux, v)

		Convey("When I add routes to an invalid view", func() {
			routeAdded := mvc.Route("GET /hello", func(_ context.Context) (mvc.View, *test.PageData, error) {
				return mvc.View("hello.tmpl"), &test.Model, nil
			})
			staticAdded := mvc.Static("GET /robots.txt", "static/bob.txt")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(staticAdded, ShouldBeFalse)
			})
		})
	})
}
