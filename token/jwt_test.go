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

func TestJWT(t *testing.T) {
	cfg := test.NewToken("jwt", "secrets/none")
	kid := token.NewKID(cfg)
	ec := test.NewEd25519()
	signer, _ := ed25519.NewSigner(ec)
	verifier, _ := ed25519.NewVerifier(ec)
	jwt := token.NewJWT(kid, signer, verifier, &id.UUID{})

	Convey("When I generate a JWT token", t, func() {
		token, err := jwt.Generate("test", "test", "test", time.Hour)
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			sub, err := jwt.Verify(token, "test", "test")
			So(err, ShouldBeNil)

			So(sub, ShouldEqual, "test")
		})

		Convey("Then I should have an error due to invalid aud", func() {
			_, err := jwt.Verify(token, "bob", "test")
			So(err, ShouldBeError)
		})

		Convey("Then I should have an error due to invalid iss", func() {
			_, err := jwt.Verify(token, "test", "bob")
			So(err, ShouldBeError)
		})
	})

	tokens := []string{
		"invalid",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	for _, token := range tokens {
		Convey("When I verify an invalid token", t, func() {
			_, err := jwt.Verify(token, "test", "test")

			Convey("Then I should have a error", func() {
				So(err, ShouldBeError)
			})
		})
	}
}
