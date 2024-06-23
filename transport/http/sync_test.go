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
	sh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/test"
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

type Errorer struct{}

func (*Errorer) Error(ctx context.Context, err error) *Response {
	return &Response{Meta: meta.CamelStrings(ctx, ""), Error: &Error{Message: err.Error()}}
}

func (*Errorer) Status(error) int {
	return http.StatusInternalServerError
}

//nolint:dupl,funlen
func TestSync(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := sh.NewServeMux(sh.StandardMux, test.RuntimeMux, sh.NewStandardServeMux())
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
			h := sh.NewHandler[Request, Response](mux, test.Marshaller, &Errorer{})

			err = h.Handle("POST", "/hello", func(_ context.Context, r *Request) (*Response, error) {
				s := "Hello " + r.Name

				return &Response{Greeting: &s}, nil
			})
			So(err, ShouldBeNil)

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
					So(resp.StatusCode, ShouldEqual, 200)
				})

				lc.RequireStop()
			})
		})
	}
}

func TestBadSync(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			mux := sh.NewServeMux(sh.StandardMux, test.RuntimeMux, sh.NewStandardServeMux())
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
			h := sh.NewHandler[Request, Response](mux, test.Marshaller, &Errorer{})

			err = h.Handle("POST", "/hello", func(_ context.Context, _ *Request) (*Response, error) {
				return nil, errors.New("ohh no")
			})
			So(err, ShouldBeNil)

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

//nolint:dupl
func TestAllowedSync(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			verifier := test.NewVerifier("test")
			mux := sh.NewServeMux(sh.StandardMux, test.RuntimeMux, sh.NewStandardServeMux())
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux, Verifier: verifier, VerifyAuth: true}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Generator: test.NewGenerator("test", nil)}
			h := sh.NewHandler[Request, Response](mux, test.Marshaller, &Errorer{})

			err := h.Handle("POST", "/hello", func(_ context.Context, r *Request) (*Response, error) {
				s := "Hello " + r.Name

				return &Response{Greeting: &s}, nil
			})
			So(err, ShouldBeNil)

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
			verifier := test.NewVerifier("test")
			mux := sh.NewServeMux(sh.StandardMux, test.RuntimeMux, sh.NewStandardServeMux())
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			cfg := test.NewInsecureTransportConfig()
			tc := test.NewOTLPTracerConfig()
			m := test.NewOTLPMeter(lc)

			s := &test.Server{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Mux: mux, Verifier: verifier, VerifyAuth: true}
			s.Register()

			cl := &test.Client{Lifecycle: lc, Logger: logger, Tracer: tc, Transport: cfg, Meter: m, Generator: test.NewGenerator("bob", nil)}
			h := sh.NewHandler[Request, Response](mux, test.Marshaller, &Errorer{})

			err := h.Handle("POST", "/hello", func(_ context.Context, r *Request) (*Response, error) {
				s := "Hello " + r.Name

				return &Response{Greeting: &s}, nil
			})
			So(err, ShouldBeNil)

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
