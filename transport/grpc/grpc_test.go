package grpc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/logger/zap"
	smeta "github.com/alexfalkowski/go-service/security/meta"
	"github.com/alexfalkowski/go-service/test"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/grpc/ratelimit"
	"github.com/alexfalkowski/go-service/transport/grpc/security/jwt"
	tmeta "github.com/alexfalkowski/go-service/transport/meta"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		serverParams := tgrpc.ServerParams{Config: cfg, Logger: logger}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
			defer cancel()

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			resp, err := client.SayHello(ctx, req)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			_, err = client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			_, err = client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			_, err = client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

// nolint:goerr113
func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("bob", errors.New("token error")))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			_, err = client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		serverParams := tgrpc.ServerParams{Config: cfg, Logger: logger}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Minute))
			defer cancel()

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			resp, err := stream.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			resp, err := stream.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestEmptyAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			_, err = client.SayStreamHello(ctx)

			Convey("Then I should have an auth error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestMissingClientAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

// nolint:goerr113
func TestTokenErrorAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("", errors.New("token error")))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			_, err = client.SayStreamHello(ctx)

			Convey("Then I should have an error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			lc.RequireStop()
		})
	})
}

func TestRateLimitUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{ratelimit.UnaryServerInterceptor(&cfg.RateLimit, tmeta.UserAgent)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			client.SayHello(ctx, req) // nolint:errcheck
			_, err = client.SayHello(ctx, req)

			Convey("Then I should have an invalid reply with rate limit", func() {
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})

			lc.RequireStop()
		})
	})
}

func TestAuthRateLimitUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary: []grpc.UnaryServerInterceptor{
				jwt.UnaryServerInterceptor(verifier),
				ratelimit.UnaryServerInterceptor(&cfg.RateLimit, smeta.AuthorizedParty),
			},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for an authenticated greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("test", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			client.SayHello(ctx, req) // nolint:errcheck
			_, err = client.SayHello(ctx, req)

			Convey("Then I should have an invalid reply with rate limit", func() {
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})

			lc.RequireStop()
		})
	})
}

func TestSuccessRateLimitStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Stream: []grpc.StreamServerInterceptor{ratelimit.StreamServerInterceptor(&cfg.RateLimit, tmeta.UserAgent)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a valid reply", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestFailedRateLimitStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		rcfg := &ratelimit.Config{
			Every: 1 * time.Microsecond,
			Burst: 0,
		}
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Stream: []grpc.StreamServerInterceptor{ratelimit.StreamServerInterceptor(rcfg, tmeta.UserAgent)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(false))

		lc.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)

			stream, err := client.SayStreamHello(ctx)
			So(err, ShouldBeNil)

			err = stream.Send(&test.HelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have an invalid reply with rate limit", func() {
				So(status.Code(err), ShouldEqual, codes.ResourceExhausted)
			})

			lc.RequireStop()
		})
	})
}

func TestBreakerUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		cfg := test.NewGRPCConfig()
		verifier := test.NewVerifier("test")
		serverParams := tgrpc.ServerParams{
			Config: cfg,
			Logger: logger,
			Unary:  []grpc.UnaryServerInterceptor{jwt.UnaryServerInterceptor(verifier)},
			Stream: []grpc.StreamServerInterceptor{jwt.StreamServerInterceptor(verifier)},
		}
		gs := tgrpc.NewServer(lc, test.NewShutdowner(), serverParams)

		test.RegisterGreeterServer(gs, test.NewServer(true))

		lc.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()
			clientParams := &tgrpc.ClientParams{
				Host:   fmt.Sprintf("127.0.0.1:%s", cfg.Port),
				Config: cfg,
				Logger: logger,
			}
			clientOpts := []grpc.DialOption{
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithPerRPCCredentials(jwt.NewPerRPCCredentials(test.NewGenerator("bob", nil))),
			}

			conn, err := tgrpc.NewClient(ctx, clientParams, clientOpts...)
			So(err, ShouldBeNil)

			defer conn.Close()

			client := test.NewGreeterClient(conn)
			req := &test.HelloRequest{Name: "test"}

			for i := 0; i < 10; i++ {
				_, err = client.SayHello(ctx, req)
			}

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unavailable)
			})

			lc.RequireStop()
		})
	})
}
