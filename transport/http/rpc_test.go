//nolint:varnamelen
package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

//nolint:gochecknoinits
func init() {
	tm.RegisterKeys()
}

func TestRPCNoContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Compression: true}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.NoContent)

			lc.RequireStart()

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(cl.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)

				_, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestRPCWithContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Compression: true}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.SuccessSayHello)

			lc.RequireStart()

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(cl.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)

				resp, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.Greeting, ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestProtobufRPC(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.ProtobufSayHello)

			lc.RequireStart()

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rpc.NewClient[v1.SayHelloRequest, v1.SayHelloResponse](url, rpc.WithClientContentType("application/"+mt))

				resp, err := client.Invoke(context.Background(), &v1.SayHelloRequest{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.GetMessage(), ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestBadUnmarshalRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.SuccessSayHello)

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				d := []byte("a bad payload")

				req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), bytes.NewReader(d))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(strings.TrimSpace(string(body)), ShouldNotBeBlank)
					So(resp.StatusCode, ShouldEqual, 400)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestErrorRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "1s", 100))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.ErrorSayHello)

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				enc := test.Encoder.Get(mt)

				b := test.Pool.Get()
				defer test.Pool.Put(b)

				err := enc.Encode(b, test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("http://%s/hello", cfg.HTTP.Address), b)
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(strings.TrimSpace(string(body)), ShouldEqual, "rpc: ohh no")
					So(resp.StatusCode, ShouldEqual, 503)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestAllowedRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			verifier := test.NewVerifier("test")
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux, Verifier: verifier, VerifyAuth: true}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Generator: test.NewGenerator("test", nil)}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.SuccessSayHello)

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(cl.NewHTTP().Transport))

				resp, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.Greeting, ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestDisallowedRPC(t *testing.T) {
	for _, mt := range []string{"application/json", "application/yaml", "application/yml", "application/toml", "application/gob", "test"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			verifier := test.NewVerifier("test")
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux, Verifier: verifier, VerifyAuth: true}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Generator: test.NewGenerator("bob", nil)}

			rpc.Register(mux, test.Content, test.Pool)
			rpc.Route("/hello", test.SuccessSayHello)

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
				url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType(mt),
					rpc.WithClientRoundTripper(cl.NewHTTP().Transport))

				_, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})

				Convey("Then I should have an error", func() {
					So(status.IsError(err), ShouldBeTrue)
					So(status.Code(err), ShouldEqual, http.StatusUnauthorized)
					So(err.Error(), ShouldContainSubstring, "token: invalid match")
				})

				lc.RequireStop()
			})
		})
	}
}

func BenchmarkRPC(b *testing.B) {
	b.ReportAllocs()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(b)
	logger := zap.NewNop()

	cfg := test.NewInsecureTransportConfig()
	tc := test.NewOTLPTracerConfig()
	m := test.NewOTLPMeter(lc)

	s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux}
	s.Register()

	cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}
	t := cl.NewHTTP().Transport

	rpc.Register(mux, test.Content, test.Pool)
	rpc.Route("/hello", test.SuccessSayHello)

	url := fmt.Sprintf("http://%s/hello", cfg.HTTP.Address)

	lc.RequireStart()
	b.ResetTimer()

	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		b.Run(mt, func(b *testing.B) {
			for range b.N {
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(t))

				_, _ = client.Invoke(context.Background(), &test.Request{Name: "Bob"})
			}
		})
	}

	b.StopTimer()
	lc.RequireStop()
}
