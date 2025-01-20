package http_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"golang.org/x/net/html"
)

func TestRouteSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression())
		world.Start()

		world.Router.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
			return mvc.View("hello.tmpl"), &test.Model
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.Request(context.Background(), "http", http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.Stop()
		})
	})
}

func TestRouteMissingView(t *testing.T) {
	Convey("Given I have setup a route with an missisng view", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Start()

		world.Router.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
			return mvc.View("none.tmpl"), &test.Model
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.Request(context.Background(), "http", http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have no html", func() {
				So(body, ShouldBeEmpty)
				So(res.StatusCode, ShouldEqual, 500)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")
			})

			world.Stop()
		})
	})
}

func TestRouteError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport))
		world.Start()

		world.Router.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
			return mvc.View("error.tmpl"), status.Error(http.StatusServiceUnavailable, "ohh no")
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.Request(context.Background(), "http", http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.Stop()
		})
	})
}

func TestStaticSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Start()

		world.Router.Static("GET /robots.txt", "static/robots.txt")

		Convey("When I query for robots", func() {
			header := http.Header{}
			header.Set("Content-Type", "text/html")

			res, body, err := world.Request(context.Background(), "http", http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get("Content-Type"), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.Stop()
		})
	})
}

func TestStaticError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Start()

		world.Router.Static("GET /robots.txt", "static/bob.txt")

		Convey("When I query for hello", func() {
			header := http.Header{}

			res, _, err := world.Request(context.Background(), "http", http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(res.StatusCode, ShouldEqual, 500)
			})

			world.Stop()
		})
	})
}

func TestMissingViews(t *testing.T) {
	Convey("Given I have views with missing FS", t, func() {
		mux := http.NewServeMux()

		v := mvc.NewViews(mvc.ViewsParams{Patterns: mvc.Patterns{"views/*.tmpl"}})
		r := mvc.NewRouter(mux, v)

		Convey("When I add routes to an invalid view", func() {
			routeAdded := r.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
				return mvc.View("hello.tmpl"), &test.Model
			})
			staticAdded := r.Static("GET /robots.txt", "static/bob.txt")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(staticAdded, ShouldBeFalse)
			})
		})
	})

	Convey("Given I have views with missing patterns", t, func() {
		mux := http.NewServeMux()
		v := mvc.NewViews(mvc.ViewsParams{FS: &test.Views})
		r := mvc.NewRouter(mux, v)

		Convey("When I add routes to an invalid view", func() {
			routeAdded := r.Route("GET /hello", func(_ context.Context) (mvc.View, mvc.Model) {
				return mvc.View("hello.tmpl"), &test.Model
			})
			staticAdded := r.Static("GET /robots.txt", "static/bob.txt")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(staticAdded, ShouldBeFalse)
			})
		})
	})
}
