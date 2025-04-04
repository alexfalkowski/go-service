package token_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidPaseto(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(ec)
	verifier, _ := ed25519.NewVerifier(ec)
	paseto := token.NewPaseto(signer, verifier, &id.UUID{})

	Convey("When I generate a paseto token", t, func() {
		token, err := paseto.Generate("test", "test", "test", time.Hour)
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			sub, err := paseto.Verify(token, "test", "test")
			So(err, ShouldBeNil)

			So(sub, ShouldEqual, "test")
		})
	})
}

func TestInvalidPaseto(t *testing.T) {
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(ec)
	verifier, _ := ed25519.NewVerifier(ec)

	Convey("When I generate a paseto token", t, func() {
		paseto := token.NewPaseto(signer, verifier, &id.UUID{})

		token, err := paseto.Generate("test", "test", "test", time.Hour)
		So(err, ShouldBeNil)

		Convey("Then I should have an error due to invalid aud", func() {
			_, err := paseto.Verify(token, "bob", "test")
			So(err, ShouldBeError)
		})

		Convey("Then I should have an error due to invalid iss", func() {
			_, err := paseto.Verify(token, "test", "bob")
			So(err, ShouldBeError)
		})
	})

	tokens := []string{"invalid"}

	for _, tkn := range tokens {
		paseto := token.NewPaseto(signer, verifier, &id.UUID{})

		Convey("When I verify an invalid token", t, func() {
			_, err := paseto.Verify(tkn, "test", "test")

			Convey("Then I should have a error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have paseto with an erroneous settings", t, func() {
		Convey("When I generate a token", func() {
			paseto := token.NewPaseto(&ed25519.Signer{}, verifier, &id.UUID{})

			_, err := paseto.Generate("test", "test", "test", time.Hour)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I verify a token", func() {
			paseto := token.NewPaseto(signer, &ed25519.Verifier{}, &id.UUID{})

			_, err := paseto.Verify("", "bob", "test")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
