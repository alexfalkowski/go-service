package http_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	th.Register(test.FS)
}

func TestRestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestNoContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				err := world.Rest.Do(t.Context(), method, url, rest.NoOptions)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})

				world.RequireStop()
			})
		})
	}
}

//nolint:dupl
func TestRestRequestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				req := &test.Request{Name: "test"}
				opts := &rest.Options{
					ContentType: mime.JSONMediaType,
					Request:     req,
				}
				err := world.Rest.Do(t.Context(), method, url, opts)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestError(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP(), test.WithWorldLoggerConfig("tilt"))
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestError)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				err := world.Rest.Do(t.Context(), method, url, rest.NoOptions)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

//nolint:dupl
func TestRestRequestError(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				req := &test.Request{Name: "test"}
				opts := &rest.Options{
					ContentType: mime.JSONMediaType,
					Request:     req,
				}
				err := world.Rest.Do(t.Context(), method, url, opts)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterHandlers("/hello", test.RestContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				resp := &test.Response{}
				opts := &rest.Options{
					Response: resp,
				}
				err := world.Rest.Do(t.Context(), method, url, opts)

				Convey("Then I should have a response", func() {
					So(err, ShouldBeNil)
					So(resp.Greeting, ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
				req := &test.Request{Name: "test"}
				resp := &test.Response{}
				opts := &rest.Options{
					ContentType: mime.JSONMediaType,
					Request:     req,
					Response:    resp,
				}
				err := world.Rest.Do(t.Context(), method, url, opts)

				Convey("Then I should have a response", func() {
					So(err, ShouldBeNil)
					So(resp.Greeting, ShouldEqual, "Hello test")
				})

				world.RequireStop()
			})
		})
	}
}

func TestRestInvalidStatusCode(t *testing.T) {
	Convey("Given I have all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		test.RegisterHandlers("/hello", test.RestInvalidStatusCode)

		Convey("When I send data", func() {
			url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
			err := world.Rest.Get(t.Context(), url, rest.NoOptions)

			Convey("Then I should have a get error", func() {
				So(err, ShouldBeError)
			})

			err = world.Rest.Delete(t.Context(), url, rest.NoOptions)

			Convey("Then I should have a delete error", func() {
				So(err, ShouldBeError)
			})
		})

		test.RegisterRequestHandlers("/hello", test.RestRequestInvalidStatusCode)

		Convey("When I send request data", func() {
			url := fmt.Sprintf("http://%s/hello", world.InsecureServerHost())
			req := &test.Request{}
			opts := &rest.Options{Request: req}

			err := world.Rest.Post(t.Context(), url, opts)

			Convey("Then I should have a post error", func() {
				So(err, ShouldBeError)
			})

			err = world.Rest.Put(t.Context(), url, opts)

			Convey("Then I should have a put error", func() {
				So(err, ShouldBeError)
			})

			err = world.Rest.Patch(t.Context(), url, opts)

			Convey("Then I should have a patch error", func() {
				So(err, ShouldBeError)
			})
		})

		world.RequireStop()
	})
}
