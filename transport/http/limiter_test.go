package http_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	meta.RegisterKeys()
}

func TestGet(t *testing.T) {
	Convey("Given I have all the servers", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		tc := test.NewOTLPTracerConfig()
		m := test.NewOTLPMeter(lc)

		s := &test.Server{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			Limiter: l, Key: k, Mux: mux,
		}
		s.Register()

		cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

		mux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("hello!"))
		})

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

			client.Do(req)
			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid greet", func() {
				So(actual, ShouldEqual, "hello!")
			})

			lc.RequireStop()
		})
	})
}

func TestLimiter(t *testing.T) {
	for _, f := range []string{"user-agent", "ip"} {
		Convey("Given I have a all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig(f, "1s", 0))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{
				Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
				Limiter: l, Key: k, Mux: mux,
			}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			mux.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("hello!"))
			})

			lc.RequireStart()

			Convey("When I query for a greet", func() {
				client := cl.NewHTTP()

				req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
				So(err, ShouldBeNil)

				client.Do(req)
				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				Convey("Then I should have been rate limited", func() {
					So(resp.StatusCode, ShouldEqual, http.StatusTooManyRequests)
					So(resp.Header.Get("RateLimit"), ShouldNotBeBlank)
				})

				lc.RequireStop()
			})
		})
	}
}
