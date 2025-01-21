package http_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestTokenAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "key"} {
		Convey("Given I have a all the servers", t, func() {
			kid, _ := token.NewKID(rand.NewGenerator(rand.NewReader()))
			a, _ := ed25519.NewSigner(test.NewEd25519())
			jwt := token.NewJWT(kid, a)
			pas := token.NewPaseto(a)
			token := token.NewToken(test.NewToken(kind), jwt, pas)

			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(token, token))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I query for an authenticated greet", func() {
				header := http.Header{}
				header.Set("Content-Type", "application/json")
				header.Set("Request-Id", "test")
				header.Set("X-Forwarded-For", "127.0.0.1")
				header.Set("Geolocation", "geo:47,11")

				res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
				So(err, ShouldBeNil)

				Convey("Then I should have a valid reply", func() {
					So(res.StatusCode, ShouldEqual, 200)
					So(body, ShouldNotBeBlank)
				})

				world.RequireStop()
			})
		})
	}
}

func TestValidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for an authenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")
			header.Set("X-Forwarded-For", "127.0.0.1")

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a valid reply", func() {
				So(res.StatusCode, ShouldEqual, 200)
				So(body, ShouldNotBeBlank)
			})

			world.RequireStop()
		})
	})
}

func TestInvalidAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(res.StatusCode, ShouldEqual, 401)
				So(body, ShouldContainSubstring, `token: invalid match`)
			})

			world.RequireStop()
		})
	})
}

func TestAuthUnaryWithAppend(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")
			header.Set("Authorization", "What Invalid")

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a reply", func() {
				So(res.StatusCode, ShouldEqual, 200)
				So(body, ShouldNotBeBlank)
			})

			world.RequireStop()
		})
	})
}

func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(res.StatusCode, ShouldEqual, 401)
				So(body, ShouldContainSubstring, "invalid match")
			})

			world.RequireStop()
		})
	})
}

func TestEmptyAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("", nil), test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")

			_, _, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization is invalid")
			})

			world.RequireStop()
		})
	})
}

func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a unauthenticated reply", func() {
				So(res.StatusCode, ShouldEqual, 401)
				So(body, ShouldContainSubstring, "invalid match")
			})

			world.RequireStop()
		})
	})
}

func TestTokenErrorAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("", test.ErrGenerate), test.NewVerifier("test")))
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a greet that will generate a token error", func() {
			header := http.Header{}
			header.Set("Content-Type", "application/json")
			header.Set("Request-Id", "test")

			_, _, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "token error")
			})

			world.RequireStop()
		})
	})
}

func TestBreakerAuthUnary(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		world := test.NewWorld(t,
			test.WithWorldTelemetry("otlp"),
			test.WithWorldToken(test.NewGenerator("", test.ErrGenerate), test.NewVerifier("test")),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			ctx := context.Background()

			var err error

			for range 10 {
				header := http.Header{}
				header.Set("Content-Type", "application/json")
				header.Set("Request-Id", "test")

				_, _, err = world.ResponseWithBody(ctx, "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			}

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		world.RequireStop()
	})
}
