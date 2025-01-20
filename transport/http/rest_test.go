//nolint:varnamelen
package http_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestRestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest())
			world.Start()

			test.RegisterHandlers("/hello", test.RestNoContent)

			Convey("When I send data", func() {
				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)

				res, err := world.Rest.R().Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(rest.Error(res), ShouldBeNil)
					So(status.Code(err), ShouldEqual, http.StatusOK)
				})

				world.Stop()
			})
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest())
			world.Start()

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			Convey("When I send data", func() {
				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)
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

				world.Stop()
			})
		})
	}
}

func TestRestError(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest())
			world.Start()

			test.RegisterHandlers("/hello", test.RestError)

			Convey("When I send data", func() {
				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)

				res, err := world.Rest.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have no error", func() {
					So(rest.Error(res), ShouldBeError)
				})

				world.Stop()
			})
		})
	}
}

func TestRestRequestError(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest())
			world.Start()

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			Convey("When I send data", func() {
				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)
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

				world.Stop()
			})
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
			world.Start()

			test.RegisterHandlers("/hello", test.RestContent)

			Convey("When I send data", func() {
				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)

				resp, err := world.Rest.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(resp, ShouldNotBeNil)
				})

				world.Stop()
			})
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
			world.Start()

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			Convey("When I send data", func() {
				var (
					b    bytes.Buffer
					resp test.Response
				)

				addr := world.Server.Transport.HTTP.Address
				url := fmt.Sprintf("http://%s/hello", addr)
				enc := json.NewEncoder()
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}

				err := enc.Encode(&b, req)
				So(err, ShouldBeNil)

				res, err := world.Rest.R().SetHeaders(headers).SetBody(b.Bytes()).Execute(v, url)
				So(err, ShouldBeNil)

				b.Reset()
				b.Write(res.Body())

				err = enc.Decode(&b, &resp)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(res, ShouldNotBeNil)
					So(resp.Greeting, ShouldEqual, "Hello test")
				})

				world.Stop()
			})
		})
	}
}
