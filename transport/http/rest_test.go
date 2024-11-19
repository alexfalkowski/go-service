package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tm.RegisterKeys()
}

func TestRest(t *testing.T) {
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
			rest.Delete("/hello", test.Rest)
			rest.Get("/hello", test.Rest)
			rest.Post("/hello", test.Rest)
			rest.Put("/hello", test.Rest)

			lc.RequireStart()

			Convey("When I send data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := cl.NewHTTP()

				req, err := http.NewRequestWithContext(context.Background(), v, url, http.NoBody)
				So(err, ShouldBeNil)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				Convey("Then I should have response", func() {
					So(resp.StatusCode, ShouldEqual, 200)
				})

				lc.RequireStop()
			})
		})
	}
}
