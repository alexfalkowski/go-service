package proxy_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/elazarl/goproxy"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestInsecureProxy(t *testing.T) {
	Convey("Given I have a proxy", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldProxy())
		world.Register()

		world.Proxy.OnRequest().DoFunc(
			func(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				res := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "It worked!")

				return req, res
			})

		world.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Should not happen"))
		})

		world.RequireStart()

		Convey("When I make a request", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodGet, "", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a header set", func() {
				So(res.StatusCode, ShouldEqual, 200)
				So(body, ShouldEqual, "It worked!")
			})
		})

		world.RequireStop()
	})
}

func TestSecureProxy(t *testing.T) {
	Convey("Given I have a proxy", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldProxy(), test.WithWorldSecure())
		world.Register()

		world.Proxy.OnRequest().DoFunc(
			func(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				res := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "It worked!")

				return req, res
			})

		world.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Should not happen"))
		})

		world.RequireStart()

		Convey("When I make a request", func() {
			header := http.Header{}

			res, body, err := world.ResponseWithBody(context.Background(), "https", world.ServerHost(), http.MethodGet, "", header, http.NoBody)
			So(err, ShouldBeNil)

			Convey("Then I should have a header set", func() {
				So(res.StatusCode, ShouldEqual, 200)
				So(body, ShouldEqual, "It worked!")
			})
		})

		world.RequireStop()
	})
}
