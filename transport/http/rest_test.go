package http_test

import (
	"fmt"
	"net/http"
	"testing"

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
	for _, mt := range []string{"json", "yaml", "yml", "toml"} {
		for _, v := range []string{"DELETE", "GET", "POST", "PUT"} {
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
				rest.Delete("/hello", test.RestNoContent)
				rest.Get("/hello", test.RestNoContent)
				rest.Post("/hello", test.RestNoContent)
				rest.Put("/hello", test.RestNoContent)

				lc.RequireStart()

				Convey("When I send data", func() {
					url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
					client := rest.NewClient(
						rest.WithClientRoundTripper(cl.NewHTTP().Transport),
						rest.WithClientTimeout("10s"),
						rest.WithClientContentType("application/"+mt),
					)

					resp, err := client.R().Execute(v, url)

					Convey("Then I should have no error", func() {
						So(err, ShouldBeNil)
						So(resp.Header().Get(content.TypeKey), ShouldEqual, "application/"+mt)
					})

					lc.RequireStop()
				})
			})
		}
	}
}

func TestRestWithContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml"} {
		for _, v := range []string{"DELETE", "GET", "POST", "PUT"} {
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
				rest.Delete("/hello", test.RestContent)
				rest.Get("/hello", test.RestContent)
				rest.Post("/hello", test.RestContent)
				rest.Put("/hello", test.RestContent)

				lc.RequireStart()

				Convey("When I send data", func() {
					url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
					client := rest.NewClient(
						rest.WithClientContentType("application/" + mt),
					)

					resp, err := client.R().Execute(v, url)
					So(err, ShouldBeNil)

					Convey("Then I should have a response", func() {
						So(resp.Body(), ShouldNotBeEmpty)
						So(resp.Header().Get(content.TypeKey), ShouldEqual, "application/"+mt)
					})

					lc.RequireStop()
				})
			})
		}
	}
}
