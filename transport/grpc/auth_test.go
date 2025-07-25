package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("bob", test.ErrGenerate), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("", nil), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(t.Context(), req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet", func() {
			ctx := t.Context()
			ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-for", "127.0.0.1")
			ctx = metadata.AppendToOutgoingContext(ctx, "geolocation", "geo:47,11")

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestAuthUnaryWithAppend(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			ctx := t.Context()
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "What Invalid")

			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have a reply", func() {
				So(status.Code(err), ShouldEqual, codes.OK)
			})

			world.RequireStop()
		})
	})
}

func TestValidAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		Convey("Given I have a gRPC server", t, func() {
			cfg := test.NewToken(kind)
			ec := test.NewEd25519()
			signer, _ := ed25519.NewSigner(test.PEM, ec)
			verifier, _ := ed25519.NewVerifier(test.PEM, ec)
			gen := &id.UUID{}
			params := token.TokenParams{
				Config: cfg,
				Name:   test.Name,
				JWT: jwt.NewToken(jwt.TokenParams{
					Config:    cfg.JWT,
					Signer:    signer,
					Verifier:  verifier,
					Generator: gen,
				}),
				Paseto: paseto.NewToken(paseto.TokenParams{
					Config:    cfg.Paseto,
					Signer:    signer,
					Verifier:  verifier,
					Generator: gen,
				}),
				SSH: ssh.NewToken(test.FS, cfg.SSH),
			}
			tkn := token.NewToken(params)

			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldGRPC())
			world.Register()
			world.RequireStart()

			Convey("When I query for an authenticated greet", func() {
				conn := world.NewGRPC()
				defer conn.Close()

				client := v1.NewGreeterServiceClient(conn)
				req := &v1.SayHelloRequest{Name: "test"}

				resp, err := client.SayHello(t.Context(), req)
				So(err, ShouldBeNil)

				Convey("Then I should have a valid reply", func() {
					So(resp.GetMessage(), ShouldEqual, "Hello test")
				})

				world.RequireStop()
			})
		})
	}
}

func TestBreakerAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
			test.WithWorldCompression(),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			var err error
			for range 10 {
				_, err = client.SayHello(t.Context(), req)
			}

			Convey("Then I should have a unavailable reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unavailable)
			})
		})

		world.RequireStop()
	})
}

func TestValidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(t.Context())
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			resp, err := stream.Recv()
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(resp.GetMessage(), ShouldEqual, "Hello test")
			})

			world.RequireStop()
		})
	})
}

func TestInvalidAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(t.Context())
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestEmptyAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("", nil), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			_, err := client.SayStreamHello(t.Context())

			Convey("Then I should have an auth error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestMissingClientAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(nil, test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			stream, err := client.SayStreamHello(t.Context())
			So(err, ShouldBeNil)

			err = stream.Send(&v1.SayStreamHelloRequest{Name: "test"})
			So(err, ShouldBeNil)

			_, err = stream.Recv()

			Convey("Then I should have a unauthenticated reply", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}

func TestTokenErrorAuthStream(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("", test.ErrGenerate), test.NewVerifier("test")),
			test.WithWorldGRPC(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a greet that will generate a token error", func() {
			conn := world.NewGRPC()
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)

			_, err := client.SayStreamHello(t.Context())

			Convey("Then I should have an error", func() {
				So(status.Code(err), ShouldEqual, codes.Unauthenticated)
			})

			world.RequireStop()
		})
	})
}
