package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerate(t *testing.T) {
	for _, kind := range []string{"opaque", "jwt", "paseto", "none"} {
		cfg := test.NewToken(kind, "secrets/opaque")
		kid := token.NewKID(cfg)
		signer, _ := ed25519.NewSigner(test.NewEd25519())
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
			JWT:    token.NewJWT(kid, signer, id.Default),
			Paseto: token.NewPaseto(signer, id.Default),
		}
		token := token.NewToken(params)

		Convey("When I try to generate", t, func() {
			_, _, err := token.Generate(t.Context())

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	}
}

func TestVerify(t *testing.T) {
	for _, kind := range []string{"opaque", "jwt", "paseto", "none"} {
		cfg := test.NewToken(kind, "secrets/opaque")
		kid := token.NewKID(cfg)
		signer, _ := ed25519.NewSigner(test.NewEd25519())
		params := token.Params{
			Config: cfg,
			Name:   test.Name,
			Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
			JWT:    token.NewJWT(kid, signer, id.Default),
			Paseto: token.NewPaseto(signer, id.Default),
		}
		token := token.NewToken(params)

		Convey("Given I generate a token", t, func() {
			_, tkn, err := token.Generate(t.Context())
			So(err, ShouldBeNil)

			Convey("When I try to verify", func() {
				_, err := token.Verify(t.Context(), tkn)

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
	signer, _ := ed25519.NewSigner(test.NewEd25519())
	params := token.Params{
		Config: cfg,
		Name:   test.Name,
		Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
		JWT:    token.NewJWT(kid, signer, id.Default),
		Paseto: token.NewPaseto(signer, id.Default),
	}
	token := token.NewToken(params)

	Convey("When I generate a token", t, func() {
		_, _, err := token.Generate(t.Context())

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I verify a token", t, func() {
		_, err := token.Verify(t.Context(), nil)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}

func TestVerifyWithNoConfig(t *testing.T) {
	params := token.Params{Name: test.Name}
	token := token.NewToken(params)

	Convey("Given I generate a token", t, func() {
		_, tkn, err := token.Generate(t.Context())
		So(err, ShouldBeNil)

		Convey("When I try to verify", func() {
			_, err := token.Verify(t.Context(), tkn)

			Convey("Then I should have no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestVerifyWithMissingToken(t *testing.T) {
	cfg := test.NewToken("opaque", "secrets/opaque")
	kid := token.NewKID(cfg)
	signer, _ := ed25519.NewSigner(test.NewEd25519())
	params := token.Params{
		Config: cfg,
		Name:   test.Name,
		Opaque: token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader())),
		JWT:    token.NewJWT(kid, signer, id.Default),
		Paseto: token.NewPaseto(signer, id.Default),
	}
	token := token.NewToken(params)

	Convey("When I verify a token", t, func() {
		_, err := token.Verify(t.Context(), nil)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
