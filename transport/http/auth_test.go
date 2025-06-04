package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenAuthUnary(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto", "ssh"} {
		Convey("Given I have a all the servers", t, func() {
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

			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(tkn, tkn), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I query for an authenticated greet", func() {
				header := http.Header{}
				header.Set(content.TypeKey, mime.JSONMediaType)
				header.Set("Request-Id", "test")
				header.Set("X-Forwarded-For", "127.0.0.1")
				header.Set("Geolocation", "geo:47,11")

				res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("test", nil), test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for an authenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")
			header.Set("X-Forwarded-For", "127.0.0.1")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("bob", nil), test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")
			header.Set("Authorization", "What Invalid")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			So(err, ShouldBeNil)

			Convey("Then I should have a reply", func() {
				So(res.StatusCode, ShouldEqual, 200)
				So(body, ShouldNotBeBlank)
			})

			world.RequireStop()
		})
	})
}

//nolint:dupl
func TestMissingAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("", nil), test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			_, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))

			Convey("Then I should have an auth error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "authorization is invalid")
			})

			world.RequireStop()
		})
	})
}

//nolint:dupl
func TestMissingClientAuthUnary(t *testing.T) {
	Convey("Given I have a all the servers", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a unauthenticated greet", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
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
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(test.NewGenerator("", test.ErrGenerate), test.NewVerifier("test")), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		rpc.Route("/hello", test.SuccessSayHello)

		Convey("When I query for a greet that will generate a token error", func() {
			header := http.Header{}
			header.Set(content.TypeKey, mime.JSONMediaType)
			header.Set("Request-Id", "test")

			_, _, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, "token: generation issue")
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
			test.WithWorldHTTP(),
		)
		world.Register()
		world.RequireStart()

		Convey("When I query for a unauthenticated greet multiple times", func() {
			var err error

			for range 10 {
				header := http.Header{}
				header.Set(content.TypeKey, mime.JSONMediaType)
				header.Set("Request-Id", "test")

				_, _, err = world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString(`{"name":"test"}`))
			}

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		world.RequireStop()
	})
}
