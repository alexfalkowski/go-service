package token_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidPaseto(t *testing.T) {
	a, _ := ed25519.NewSigner(test.NewEd25519())
	paseto := token.NewPaseto(a, id.Default)

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
	Convey("When I generate a paseto token", t, func() {
		a, _ := ed25519.NewSigner(test.NewEd25519())
		paseto := token.NewPaseto(a, id.Default)

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
		a, _ := ed25519.NewSigner(test.NewEd25519())
		paseto := token.NewPaseto(a, id.Default)

		Convey("When I verify an invalid token", t, func() {
			_, err := paseto.Verify(tkn, "test", "test")

			Convey("Then I should have a errror", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have paseto with an erroneous signer", t, func() {
		paseto := token.NewPaseto(&ed25519.Signer{}, id.Default)

		Convey("When I generate a token", func() {
			_, err := paseto.Generate("test", "test", "test", time.Hour)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I verify a token", func() {
			_, err := paseto.Verify("", "bob", "test")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
