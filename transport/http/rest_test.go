package http_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tm.RegisterKeys()
}

func TestRestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rest.Register(mux, test.Content)
			registerHandlers("/hello", test.RestNoContent)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient(
					rest.WithClientRoundTripper(cl.NewHTTP().Transport),
					rest.WithClientTimeout("10s"),
				)

				res, err := client.R().Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(rest.Error(res), ShouldBeNil)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rest.Register(mux, test.Content)
			registerBodyHandlers("/hello", test.RestRequestNoContent)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient(
					rest.WithClientRoundTripper(cl.NewHTTP().Transport),
					rest.WithClientTimeout("10s"),
				)
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}
				res, err := client.R().SetHeaders(headers).SetBody(req).Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(rest.Error(res), ShouldBeNil)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestError(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rest.Register(mux, test.Content)
			registerHandlers("/hello", test.RestError)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient(
					rest.WithClientRoundTripper(cl.NewHTTP().Transport),
					rest.WithClientTimeout("10s"),
				)

				res, err := client.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have no error", func() {
					So(rest.Error(res), ShouldBeError)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestRequestError(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rest.Register(mux, test.Content)
			registerBodyHandlers("/hello", test.RestRequestError)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient(
					rest.WithClientRoundTripper(cl.NewHTTP().Transport),
					rest.WithClientTimeout("10s"),
				)
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}
				res, err := client.R().SetHeaders(headers).SetBody(req).Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have no error", func() {
					So(rest.Error(res), ShouldBeError)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodDelete, http.MethodGet} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			rest.Register(mux, test.Content)
			registerHandlers("/hello", test.RestContent)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient()

				resp, err := client.R().Execute(v, url)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(resp, ShouldNotBeNil)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, v := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
			s.Register()

			rest.Register(mux, test.Content)
			registerBodyHandlers("/hello", test.RestRequestContent)

			lc.RequireStart()

			Convey("When I send data", func() {
				var (
					b    bytes.Buffer
					resp test.Response
				)

				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rest.NewClient()
				enc := json.NewEncoder()
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				req := &test.Request{Name: "test"}

				err := enc.Encode(&b, req)
				So(err, ShouldBeNil)

				res, err := client.R().SetHeaders(headers).SetBody(b.Bytes()).Execute(v, url)
				So(err, ShouldBeNil)

				b.Reset()
				b.Write(res.Body())

				err = enc.Decode(&b, &resp)
				So(err, ShouldBeNil)

				Convey("Then I should have a response", func() {
					So(res, ShouldNotBeNil)
					So(resp.Greeting, ShouldEqual, "Hello test")
				})

				lc.RequireStop()
			})
		})
	}
}

func registerHandlers[Res any](path string, h content.ResponseHandler[Res]) {
	rest.Delete(path, h)
	rest.Get(path, h)
}

func registerBodyHandlers[Req any, Res any](path string, h content.RequestResponseHandler[Req, Res]) {
	rest.Post(path, h)
	rest.Put(path, h)
	rest.Patch(path, h)
}
