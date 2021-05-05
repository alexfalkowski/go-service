package http_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgGRPC "github.com/alexfalkowski/go-service/pkg/grpc"
	tokenGRPC "github.com/alexfalkowski/go-service/pkg/grpc/security/token"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/http"
	tokenHTTP "github.com/alexfalkowski/go-service/pkg/http/security/token"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10009", HTTPPort: "10010"}

		mux := pkgHTTP.NewMux()
		pkgHTTP.Register(lc, sh, mux, cfg, logger)

		gs := pkgGRPC.NewServer(lc, sh, cfg, logger)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, opts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a greet", func() {
			client := &http.Client{Transport: pkgHTTP.NewRoundTripper(logger)}

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10010/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"message":"Hello test"}`)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl,funlen
func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10011", HTTPPort: "10012"}

		mux := pkgHTTP.NewMux()
		pkgHTTP.Register(lc, sh, mux, cfg, logger)

		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, sh, cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for an authenticated greet", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("test"), pkgHTTP.NewRoundTripper(logger))
			client := &http.Client{Transport: transport}

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10012/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

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
		sh := test.NewShutdowner()
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{GRPCPort: "10013", HTTPPort: "10014"}

		mux := pkgHTTP.NewMux()
		pkgHTTP.Register(lc, sh, mux, cfg, logger)

		verifier := test.NewVerifier("test")
		serverUnaryOpt := pkgGRPC.UnaryServerOption(logger, tokenGRPC.UnaryServerInterceptor(verifier))
		serverStreamOpt := pkgGRPC.StreamServerOption(logger, tokenGRPC.StreamServerInterceptor(verifier))
		gs := pkgGRPC.NewServer(lc, sh, cfg, logger, serverUnaryOpt, serverStreamOpt)

		test.RegisterGreeterServer(gs, test.NewServer())

		lc.RequireStart()

		ctx := context.Background()
		clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithInsecure()}

		conn, err := pkgGRPC.NewClient(ctx, fmt.Sprintf("127.0.0.1:%s", cfg.GRPCPort), logger, clientOpts...)
		So(err, ShouldBeNil)

		defer conn.Close()

		err = test.RegisterGreeterHandler(ctx, mux, conn)
		So(err, ShouldBeNil)

		Convey("When I query for a unauthenticated greet", func() {
			transport := tokenHTTP.NewRoundTripper(test.NewGenerator("bob"), pkgHTTP.NewRoundTripper(logger))
			client := &http.Client{Transport: transport}

			message := []byte(`{"name":"test"}`)
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:10014/v1/greet/hello", bytes.NewBuffer(message))
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Request-ID", "test")

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			actual := strings.TrimSpace(string(body))

			lc.RequireStop()

			Convey("Then I should have a valid reply", func() {
				So(actual, ShouldEqual, `{"code":16,"message":"could not verify token: invalid token","details":[]}`)
			})

			lc.RequireStop()
		})
	})
}
