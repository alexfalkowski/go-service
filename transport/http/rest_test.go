package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/go-resty/resty/v2"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

var methods = []string{
	resty.MethodDelete, resty.MethodGet, resty.MethodPost,
	resty.MethodPut, resty.MethodPatch, resty.MethodHead,
	resty.MethodOptions,
}

func init() {
	tm.RegisterKeys()
}

func TestRestNoContent(t *testing.T) {
	for _, v := range methods {
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

				_, err := client.R().Execute(v, url)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, v := range methods {
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

func registerHandlers(path string, h content.Handler) {
	rest.Delete(path, h)
	rest.Get(path, h)
	rest.Post(path, h)
	rest.Put(path, h)
	rest.Patch(path, h)
	rest.Head(path, h)
	rest.Options(path, h)
}
