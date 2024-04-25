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
	hl "github.com/alexfalkowski/go-service/transport/http/limiter"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"github.com/urfave/negroni/v3"
	"go.uber.org/fx/fxtest"
)

func TestGet(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		l, err := limiter.New(test.NewLimiterConfig("100-S"))
		So(err, ShouldBeNil)

		cfg := test.NewInsecureTransportConfig()
		cfg.GRPC.Enabled = false

		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, []negroni.Handler{hl.NewHandler(l, tm.UserAgent)})
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		err = test.Mux.HandlePath("GET", "/hello", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
			w.Write([]byte("hello!"))
		})
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

			req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
			So(err, ShouldBeNil)

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
	for _, f := range []limiter.KeyFunc{tm.UserAgent, tm.IPAddr} {
		Convey("Given I have a all the servers", t, func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, err := limiter.New(test.NewLimiterConfig("0-S"))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			cfg.GRPC.Enabled = false

			m := test.NewMeter(lc)
			hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, []negroni.Handler{hl.NewHandler(l, f)})
			gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

			test.RegisterTransport(lc, gs, hs)
			lc.RequireStart()

			err = test.Mux.HandlePath("GET", "/hello", func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
				w.Write([]byte("hello!"))
			})
			So(err, ShouldBeNil)

			Convey("When I query for a greet", func() {
				client := test.NewHTTPClient(lc, logger, test.NewOTLPTracerConfig(), cfg, m)

				req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), http.NoBody)
				So(err, ShouldBeNil)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				Convey("Then I should have been rate limited", func() {
					So(resp.StatusCode, ShouldEqual, http.StatusTooManyRequests)
					So(resp.Header.Get("X-Rate-Limit-Limit"), ShouldEqual, "0")
				})

				lc.RequireStop()
			})
		})
	}
}
