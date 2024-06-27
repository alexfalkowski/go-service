package http_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	nh "github.com/alexfalkowski/go-service/net/http"
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
	Error    *Error
	Greeting *string
}

type Error struct {
	Message string
}

type SuccessHandler struct{}

func (*SuccessHandler) Handle(ctx nh.Context, r *Request) (*Response, error) {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Greeting: &s}, nil
}

func (*SuccessHandler) Error(ctx nh.Context, err error) *Response {
	return &Response{Meta: meta.CamelStrings(ctx, ""), Error: &Error{Message: err.Error()}}
}

func (*SuccessHandler) Status(error) int {
	return http.StatusInternalServerError
}

type ProtobufHandler struct{}

func (*ProtobufHandler) Handle(_ nh.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func (*ProtobufHandler) Error(_ nh.Context, err error) *v1.SayHelloResponse {
	return &v1.SayHelloResponse{Message: err.Error()}
}

func (*ProtobufHandler) Status(error) int {
	return http.StatusInternalServerError
}

type ErrorHandler struct{}

func (*ErrorHandler) Handle(_ nh.Context, _ *Request) (*Response, error) {
	return nil, errors.New("ohh no")
}

func (*ErrorHandler) Error(ctx nh.Context, err error) *Response {
	return &Response{Meta: meta.CamelStrings(ctx, ""), Error: &Error{Message: err.Error()}}
}

func (*ErrorHandler) Status(error) int {
	return http.StatusInternalServerError
}

func TestSync(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := http.NewServeMux()
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "100-S"))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Compression: true}

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &SuccessHandler{})

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

				var r Response
				err = mar.Unmarshal(body, &r)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(*r.Greeting, ShouldEqual, "Hello Bob")
					So(resp.Header.Get("Content-Type"), ShouldEqual, "application/"+mt)
					So(resp.Header.Get("Service-Version"), ShouldEqual, "1.0.0")
					So(resp.StatusCode, ShouldEqual, 200)
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

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "100-S"))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &ProtobufHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				mar := test.Marshaller.Get(mt)

				d, err := mar.Marshal(&v1.SayHelloRequest{Name: "Bob"})
				So(err, ShouldBeNil)

				req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewReader(d))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				var r v1.SayHelloResponse
				err = mar.Unmarshal(body, &r)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(r.GetMessage(), ShouldEqual, "Hello Bob")
					So(resp.Header.Get("Content-Type"), ShouldEqual, "application/"+mt)
					So(resp.StatusCode, ShouldEqual, 200)
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

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "100-S"))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post data", func() {
				client := cl.NewHTTP()
				mar := test.Marshaller.Get(mt)
				d := []byte("a bad payload")

				req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/hello", cfg.HTTP.Port), bytes.NewReader(d))
				So(err, ShouldBeNil)

				req.Header.Set("Content-Type", "application/"+mt)

				resp, err := client.Do(req)
				So(err, ShouldBeNil)

				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				So(err, ShouldBeNil)

				var r Response
				err = mar.Unmarshal(body, &r)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(r.Error.Message, ShouldNotBeBlank)
					So(resp.Header.Get("Content-Type"), ShouldEqual, "application/"+mt)
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

			l, k, err := limiter.New(test.NewLimiterConfig("user-agent", "100-S"))
			So(err, ShouldBeNil)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Limiter: l, Key: k, Mux: mux}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m}

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &ErrorHandler{})

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

				var r Response
				err = mar.Unmarshal(body, &r)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(r.Error.Message, ShouldEqual, "invalid handle: ohh no")
					So(resp.Header.Get("Content-Type"), ShouldEqual, "application/"+mt)
					So(resp.StatusCode, ShouldEqual, 500)
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

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
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

				var r Response
				err = mar.Unmarshal(body, &r)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(*r.Greeting, ShouldEqual, "Hello Bob")
					So(resp.Header.Get("Content-Type"), ShouldEqual, "application/"+mt)
					So(resp.StatusCode, ShouldEqual, 200)
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

			nh.Register(mux, test.Marshaller)
			nh.Handle("/hello", &SuccessHandler{})

			lc.RequireStart()

			Convey("When I post authenticated data", func() {
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
					So(strings.TrimSpace(string(body)), ShouldEqual, "verify token: invalid token")
					So(resp.StatusCode, ShouldEqual, 401)
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
