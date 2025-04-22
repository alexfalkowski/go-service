//nolint:varnamelen
package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/net/http/status"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestNoContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				res, err := world.Rest.R().Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(rest.Error(res), ShouldBeNil)
					So(status.Code(err), ShouldEqual, http.StatusOK)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}
				res, err := world.Rest.R().SetHeaders(headers).SetBody(req).Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(rest.Error(res), ShouldBeNil)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestError(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP(), test.WithWorldLoggerConfig("tilt"))
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestError)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())

				res, err := world.Rest.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have no error", func() {
					So(rest.Error(res), ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestRequestError(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}
				res, err := world.Rest.R().SetHeaders(headers).SetBody(req).Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have no error", func() {
					So(rest.Error(res), ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())

				resp, err := world.Rest.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(resp, ShouldNotBeNil)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			Convey("When I send data", func() {
				var resp test.Response

				b := test.Pool.Get()
				defer test.Pool.Put(b)

				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				enc := json.NewEncoder()
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}

				err := enc.Encode(b, req)
				So(err, ShouldBeNil)

				res, err := world.Rest.R().SetHeaders(headers).SetBody(b.Bytes()).Execute(v, url)
				So(err, ShouldBeNil)

				b.Reset()
				b.Write(res.Body())

				err = enc.Decode(b, &resp)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(res, ShouldNotBeNil)
					So(resp.Greeting, ShouldEqual, "Hello test")
				})

				world.RequireStop()
			})
		})
	}
}
