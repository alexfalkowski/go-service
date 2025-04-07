package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"opaque", "jwt", "paseto", "none"} {
		cfg := test.NewToken(kind, "secrets/opaque")
		kid := token.NewKID(cfg)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(ec)
		verifier, _ := ed25519.NewVerifier(ec)
		gen := &id.UUID{}
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
			JWT:    token.NewJWT(kid, signer, verifier, gen),
			Paseto: token.NewPaseto(signer, verifier, gen),
		}
		tkn := token.NewToken(params)

		Convey("When I try to generate", t, func() {
			_, _, err := tkn.Generate(t.Context())

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"jwt", "paseto"} {
		cfg := test.NewToken(kind, "secrets/opaque")
		kid := token.NewKID(cfg)
		ec := test.NewEd25519()
		signer, _ := ed25519.NewSigner(ec)
		verifier, _ := ed25519.NewVerifier(ec)
		gen := &id.UUID{}
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			JWT:    token.NewJWT(kid, signer, verifier, gen),
			Paseto: token.NewPaseto(signer, verifier, gen),
		}
		tkn := token.NewToken(params)

		Convey("Given I generate a token", t, func() {
			_, gen, err := tkn.Generate(t.Context())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				ctx, err := tkn.Verify(t.Context(), gen)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(token.Subject(ctx).String(), ShouldEqual, "sub")
				})
			})
		})
	}

	for _, kind := range []string{"opaque", "ssh", "none"} {
		cfg := test.NewToken(kind, "secrets/opaque")
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
			SSH:    ssh.NewToken(cfg.SSH),
		}
		tkn := token.NewToken(params)

		Convey("Given I generate a token", t, func() {
			_, gen, err := tkn.Generate(t.Context())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				_, err := tkn.Verify(t.Context(), gen)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	}
}

func TestError(t *testing.T) {
	cfg := test.NewToken("opaque", "secrets/none")
	kid := token.NewKID(cfg)
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(ec)
	verifier, _ := ed25519.NewVerifier(ec)
	gen := &id.UUID{}
	params := token.Params{
		Config: cfg,
		Name:   test.Name,
		Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
		JWT:    token.NewJWT(kid, signer, verifier, gen),
		Paseto: token.NewPaseto(signer, verifier, gen),
	}
	tkn := token.NewToken(params)

	Convey("When I generate a token", t, func() {
		_, _, err := tkn.Generate(t.Context())

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I verify a token", t, func() {
		_, err := tkn.Verify(t.Context(), nil)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}

func TestWithNoConfig(t *testing.T) {
	Convey("When I try to create a token with no config", t, func() {
		params := token.Params{Name: test.Name}
		tkn := token.NewToken(params)

		Convey("Then I should have bo token", func() {
			So(tkn, ShouldBeNil)
		})
	})
}

func TestVerifyWithMissingToken(t *testing.T) {
	cfg := test.NewToken("opaque", "secrets/opaque")
	kid := token.NewKID(cfg)
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(ec)
	verifier, _ := ed25519.NewVerifier(ec)
	gen := &id.UUID{}
	params := token.Params{
		Config: cfg,
		Name:   test.Name,
		Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
		JWT:    token.NewJWT(kid, signer, verifier, gen),
		Paseto: token.NewPaseto(signer, verifier, gen),
	}
	tkn := token.NewToken(params)

	Convey("When I verify a token", t, func() {
		_, err := tkn.Verify(t.Context(), nil)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
