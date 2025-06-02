package paseto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/context"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	paseto := paseto.NewToken(cfg.Paseto, signer, verifier, &id.UUID{})

	Convey("When I generate a paseto token", t, func() {
		ctx := context.WithOpts(t.Context(), context.Options{
			"sub": "sub",
			"aud": "aud",
		})

		token, err := paseto.Generate(ctx)
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			ctx, err := paseto.Verify(ctx, token)
			So(err, ShouldBeNil)

			So(context.Opts(ctx).GetString("sub"), ShouldEqual, "sub")
		})
	})
}

//nolint:funlen
func TestInvalid(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)

	Convey("When I generate a paseto token with invalid aud", t, func() {
		cfg := test.NewToken("paseto")
		token := paseto.NewToken(cfg.Paseto, signer, verifier, &id.UUID{})
		ctx := context.WithOpts(t.Context(), context.Options{
			"sub": "sub",
			"aud": "test",
		})

		tkn, err := token.Generate(ctx)
		So(err, ShouldBeNil)

		token = paseto.NewToken(cfg.Paseto, signer, verifier, &id.UUID{})

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			ctx := context.WithOpts(t.Context(), context.Options{
				"aud": "aud",
			})

			_, err := token.Verify(ctx, tkn)
			So(err, ShouldBeError)
		})
	})

	Convey("When I generate a JWT token with invalid iss", t, func() {
		pcf := &paseto.Config{
			Issuer:     "test",
			Expiration: "1h",
		}
		token := paseto.NewToken(pcf, signer, verifier, &id.UUID{})
		ctx := context.WithOpts(t.Context(), context.Options{
			"sub": "sub",
			"aud": "aud",
		})

		tkn, err := token.Generate(ctx)
		So(err, ShouldBeNil)

		cfg := test.NewToken("paseto")
		token = paseto.NewToken(cfg.Paseto, signer, verifier, &id.UUID{})

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			_, err := token.Verify(ctx, tkn)
			So(err, ShouldBeError)
		})
	})

	for _, tkn := range []string{"invalid"} {
		cfg := test.NewToken("paseto")
		token := paseto.NewToken(cfg.Paseto, signer, verifier, &id.UUID{})
		ctx := context.WithOpts(t.Context(), context.Options{
			"sub": "sub",
			"aud": "aud",
		})

		Convey("When I verify an invalid token", t, func() {
			_, err := token.Verify(ctx, tkn)

			Convey("Then I should have a error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have paseto with an erroneous settings", t, func() {
		cfg := test.NewToken("paseto")
		ctx := context.WithOpts(t.Context(), context.Options{
			"sub": "sub",
			"aud": "aud",
		})

		Convey("When I generate a token", func() {
			token := paseto.NewToken(cfg.Paseto, &ed25519.Signer{}, verifier, &id.UUID{})

			_, err := token.Generate(ctx)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I verify a token", func() {
			token := paseto.NewToken(cfg.Paseto, signer, &ed25519.Verifier{}, &id.UUID{})

			_, err := token.Verify(ctx, "")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
