package paseto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	cfg := test.NewToken("paseto")
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(test.PEM, ec)
	verifier, _ := ed25519.NewVerifier(test.PEM, ec)
	paseto := paseto.NewToken(paseto.TokenParams{
		Config:    cfg.Paseto,
		Signer:    signer,
		Verifier:  verifier,
		Generator: &id.UUID{},
	})

	Convey("When I generate a paseto token", t, func() {
		token, err := paseto.Generate("hello", test.UserID.String())
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			sub, err := paseto.Verify(token, "hello")
			So(err, ShouldBeNil)

			So(sub, ShouldEqual, test.UserID.String())
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
		paseto := paseto.NewToken(paseto.TokenParams{
			Config:    cfg.Paseto,
			Signer:    signer,
			Verifier:  verifier,
			Generator: &id.UUID{},
		})

		token, err := paseto.Generate("hello", test.UserID.String())
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			_, err := paseto.Verify(token, "test")
			So(err, ShouldBeError)
		})
	})

	Convey("When I generate a JWT token with invalid iss", t, func() {
		pcf := &paseto.Config{
			Issuer:     "test",
			Expiration: "1h",
		}
		token := paseto.NewToken(paseto.TokenParams{
			Config:    pcf,
			Signer:    signer,
			Verifier:  verifier,
			Generator: &id.UUID{},
		})

		tkn, err := token.Generate("hello", test.UserID.String())
		So(err, ShouldBeNil)

		cfg := test.NewToken("paseto")
		token = paseto.NewToken(paseto.TokenParams{
			Config:    cfg.Paseto,
			Signer:    signer,
			Verifier:  verifier,
			Generator: &id.UUID{},
		})

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should have an error", func() {
			_, err := token.Verify(tkn, "hello")
			So(err, ShouldBeError)
		})
	})

	for _, tkn := range []string{"invalid"} {
		cfg := test.NewToken("paseto")
		token := paseto.NewToken(paseto.TokenParams{
			Config:    cfg.Paseto,
			Signer:    signer,
			Verifier:  verifier,
			Generator: &id.UUID{},
		})

		Convey("When I verify an invalid token", t, func() {
			_, err := token.Verify(tkn, "aud")

			Convey("Then I should have a error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have paseto with an erroneous settings", t, func() {
		cfg := test.NewToken("paseto")

		Convey("When I generate a token", func() {
			token := paseto.NewToken(paseto.TokenParams{
				Config:    cfg.Paseto,
				Signer:    &ed25519.Signer{},
				Verifier:  verifier,
				Generator: &id.UUID{},
			})

			_, err := token.Generate("hello", test.UserID.String())

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I verify a token", func() {
			token := paseto.NewToken(paseto.TokenParams{
				Config:    cfg.Paseto,
				Signer:    signer,
				Verifier:  &ed25519.Verifier{},
				Generator: &id.UUID{},
			})

			_, err := token.Verify("", "aud")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
