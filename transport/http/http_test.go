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
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tm.RegisterKeys()
}

type Request struct {
	Name string
}

type Response struct {
	Meta     meta.Map
	Greeting *string
}

type SuccessHandler struct{}

func (*SuccessHandler) Handle(ctx rpc.Context, r *Request) (*Response, error) {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Greeting: &s}, nil
}

type ProtobufHandler struct{}

func (*ProtobufHandler) Handle(_ rpc.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func (*ProtobufHandler) Error(_ rpc.Context, err error) *v1.SayHelloResponse {
	return &v1.SayHelloResponse{Message: err.Error()}
}

type ErrorHandler struct{}

func (*ErrorHandler) Handle(_ rpc.Context, _ *Request) (*Response, error) {
	return nil, rpc.Error(http.StatusServiceUnavailable, "ohh no")
}

func TestSync(t *testing.T) {
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

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Compression: true, H2C: true}

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := rpc.NewClient[Request, Response](fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), "application/"+mt, cl.NewHTTP(), test.Marshaller)

				resp, err := client.Call(context.Background(), &Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(*resp.Greeting, ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestProtobufSync(t *testing.T) {
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

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &ProtobufHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := rpc.NewClient[v1.SayHelloRequest, v1.SayHelloResponse](fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), "application/"+mt, cl.NewHTTP(), test.Marshaller)

				resp, err := client.Call(context.Background(), &v1.SayHelloRequest{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.GetMessage(), ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestBadUnmarshalSync(t *testing.T) {
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

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				d := []byte("a bad payload")

				req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewReader(d))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(strings.TrimSpace(string(body)), ShouldNotBeBlank)
					So(resp.StatusCode, ShouldEqual, 500)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestErrorSync(t *testing.T) {
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

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &ErrorHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				mar := test.Marshaller.Get(mt)

				d, err := mar.Marshal(Request{Name: "Bob"})
				So(err, ShouldBeNil)

				req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewReader(d))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(strings.TrimSpace(string(body)), ShouldEqual, "rpc handle: ohh no")
					So(resp.StatusCode, ShouldEqual, 503)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestAllowedSync(t *testing.T) {
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

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
				client := rpc.NewClient[Request, Response](fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), "application/"+mt, cl.NewHTTP(), test.Marshaller)

				resp, err := client.Call(context.Background(), &Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(*resp.Greeting, ShouldEqual, "Hello Bob")
				})

				lc.RequireStop()
			})
		})
	}
}

func TestDisallowedSync(t *testing.T) {
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

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Generator: test.NewGenerator("bob", nil)}

			rpc.Register(mux, test.Marshaller)
			rpc.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
				client := rpc.NewClient[Request, Response](fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), "application/"+mt, cl.NewHTTP(), test.Marshaller)

				_, err := client.Call(context.Background(), &Request{Name: "Bob"})

				Convey("Then I should have an error", func() {
					So(rpc.IsError(err), ShouldBeTrue)
					So(rpc.Code(err), ShouldEqual, http.StatusUnauthorized)
					So(err.Error(), ShouldContainSubstring, "verify token: invalid token")
				})

				lc.RequireStop()
			})
		})
	}
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
