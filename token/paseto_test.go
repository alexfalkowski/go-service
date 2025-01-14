package token_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestPaseto(t *testing.T) {
	a, _ := ed25519.NewSigner(test.NewEd25519())
	paseto := token.NewPaseto(a)

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

	for _, token := range tokens {
		Convey("When I verify an invalid token", t, func() {
			_, err := paseto.Verify(token, "test", "test")

			Convey("Then I should have a errror", func() {
				So(err, ShouldBeError)
			})
		})
	}
}
