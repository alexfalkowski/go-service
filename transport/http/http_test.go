package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/test"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tm.RegisterKeys()
}

func TestSecure(t *testing.T) {
	Convey("Given I a secure client", t, func() {
		mux := http.NewServeMux()
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		tc := test.NewOTLPTracerConfig()
		m := test.NewPrometheusMeter(lc)
		cfg := test.NewSecureTransportConfig()

		s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
		s.Register()

		cl := &test.Client{
			Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m,
			TLS: test.NewTLSClientConfig(),
		}

		lc.RequireStart()

		Convey("When I query github", func() {
			client := cl.NewHTTP()

			req, err := http.NewRequestWithContext(context.Background(), "GET", "https://github.com/alexfalkowski", http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have valid response", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

		lc.RequireStop()
	})
}

func BenchmarkHTTP(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	p := test.Port()

	mux.HandleFunc("GET /hello", func(_ http.ResponseWriter, _ *http.Request) {
	})

	server := &http.Server{
		Handler:           mux,
		Addr:              ":" + p,
		ReadHeaderTimeout: time.Second,
	}
	defer server.Close()

	go server.ListenAndServe()

	client := &http.Client{Transport: http.DefaultTransport}
	url := fmt.Sprintf("http://localhost:%s/hello", p)

	b.ResetTimer()

	b.Run("std", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, err := http.NewRequestWithContext(context.Background(), "GET", url, http.NoBody)
			runtime.Must(err)

			_, err = client.Do(req)
			runtime.Must(err)
		}
	})

	b.StopTimer()
}
