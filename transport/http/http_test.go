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
	"time"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	jgrpc "github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	hprometheus "github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	jhttp "github.com/alexfalkowski/go-service/transport/http/security/jwt"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), false, nil, nil)

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version))

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")
			req.Header.Set("X-Forwarded-For", "test")
			req.Header.Set("Geolocation", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

func TestDefaultClientUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewDatadogConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		hs, hport := test.NewHTTPServer(lc, logger, test.NewDatadogConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewDatadogConfig(), false, nil, nil)

		lc.RequireStart()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
		defer cancel()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := http.DefaultClient

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()
		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("test", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version), transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl,funlen
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("bob", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version), transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, `could not verify token: invalid token`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version))

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, `authorization token is not provided`)
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", nil), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version), transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req)

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization token is not provided")
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			client := test.NewHTTPClient(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version))

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			Convey("Then I should have a unauthenticated reply", func() {
				So(actual, ShouldContainSubstring, `authorization token is not provided`)
			})

			lc.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		verifier := test.NewVerifier("test")
		hs, hport := test.NewHTTPServer(lc, logger, test.NewJaegerConfig())
		_, gconfig := test.NewGRPCServer(lc, logger, test.NewJaegerConfig(), true,
			[]grpc.UnaryServerInterceptor{jgrpc.UnaryServerInterceptor(verifier)},
			[]grpc.StreamServerInterceptor{jgrpc.StreamServerInterceptor(verifier)},
		)

		lc.RequireStart()

		ctx := context.Background()

		conn, err := tgrpc.NewClient(
			tgrpc.ClientParams{Context: ctx, Host: fmt.Sprintf("127.0.0.1:%s", gconfig.Port), Version: test.Version, Config: gconfig},
			tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
			tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(),
			tgrpc.WithClientDialOption(grpc.WithBlock()),
		)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = v1.RegisterGreeterServiceHandler(ctx, hs.Mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet that will generate a token error", func() {
			transport := jhttp.NewRoundTripper(test.NewGenerator("", errors.New("token error")), http.DefaultTransport)
			client := test.NewHTTPClientWithRoundTripper(logger, tracer, test.Version, hprometheus.NewClientMetrics(lc, test.Version), transport)

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%s/v1/greet/hello", hport), bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			_, err = client.Do(req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "token error")
			})

			lc.RequireStop()
		})
	})
}
