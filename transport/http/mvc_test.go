package http_test

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/net/http/content"
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

		full, _ := mvc.NewViewPair("views/hello.tmpl")

		controller := func(_ context.Context) (*mvc.View, *test.Page, error) {
			return full, &test.Model, nil
		}

		mvc.Delete("/hello", controller)
		mvc.Get("/hello", controller)
		mvc.Post("/hello", controller)
		mvc.Put("/hello", controller)
		mvc.Patch("/hello", controller)

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set(content.TypeKey, "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get(content.TypeKey), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	})
}

func TestRoutePartialViewSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldCompression(), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		_, partial := mvc.NewViewPair("views/hello.tmpl")

		mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
			return partial, &test.Model, nil
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set(content.TypeKey, "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get(content.TypeKey), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
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

		mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
			return mvc.NewFullView("views/error.tmpl"), &test.Model, status.FromError(http.StatusServiceUnavailable, test.ErrFailed)
		})

		Convey("When I query for hello", func() {
			header := http.Header{}
			header.Set(content.TypeKey, "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "hello", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have an error", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 503)
				So(res.Header.Get(content.TypeKey), ShouldEqual, "text/html; charset=utf-8")

				_, err := html.Parse(strings.NewReader(body))
				So(err, ShouldBeNil)
			})

			world.RequireStop()
		})
	})
}

func TestStaticFileSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.StaticFile("/robots.txt", "static/robots.txt")

		Convey("When I query for robots", func() {
			header := http.Header{}
			header.Set(content.TypeKey, "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get(content.TypeKey), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.RequireStop()
		})
	})
}

func TestStaticFileError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.StaticFile("/robots.txt", "static/bob.txt")

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

func TestStaticPathValueSuccess(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.StaticPathValue("/{file}", "file", "static")

		Convey("When I query for robots", func() {
			header := http.Header{}
			header.Set(content.TypeKey, "text/html")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "robots.txt", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have valid html", func() {
				So(body, ShouldNotBeEmpty)
				So(res.StatusCode, ShouldEqual, 200)
				So(res.Header.Get(content.TypeKey), ShouldEqual, "text/plain; charset=utf-8")
			})

			world.RequireStop()
		})
	})
}

func TestStaticPathValueError(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		mvc.StaticPathValue("/{file}", "file", "static")

		Convey("When I query for hello", func() {
			header := http.Header{}

			res, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodGet, "bob.txt", header, http.NoBody)
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
		mvc.Register(mvc.RegisterParams{
			Mux:         http.NewServeMux(),
			FunctionMap: mvc.NewFunctionMap(slog.Default()),
			Layout:      test.Layout,
		})

		Convey("When I add routes to an invalid view", func() {
			routeAdded := mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
				return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
			})
			fileAdded := mvc.StaticFile("/robots.txt", "static/robots.txt")
			pathAdded := mvc.StaticPathValue("/{file}", "file", "static")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(fileAdded, ShouldBeFalse)
				So(pathAdded, ShouldBeFalse)
			})
		})
	})

	Convey("Given I have views with missing layout", t, func() {
		mvc.Register(mvc.RegisterParams{
			Mux:         http.NewServeMux(),
			FunctionMap: mvc.NewFunctionMap(slog.Default()),
			FileSystem:  test.FileSystem,
		})

		Convey("When I add routes to an invalid view", func() {
			routeAdded := mvc.Get("/hello", func(_ context.Context) (*mvc.View, *test.Page, error) {
				return mvc.NewFullView("views/hello.tmpl"), &test.Model, nil
			})
			fileAdded := mvc.StaticFile("/robots.txt", "static/robots.txt")
			pathAdded := mvc.StaticPathValue("/{file}", "file", "static")

			Convey("Then they should not be added", func() {
				So(routeAdded, ShouldBeFalse)
				So(fileAdded, ShouldBeFalse)
				So(pathAdded, ShouldBeFalse)
			})
		})
	})
}
