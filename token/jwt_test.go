package token_test

import (
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestKID(t *testing.T) {
	Convey("When I generate a KID", t, func() {
		rand := rand.NewGenerator(rand.NewReader())
		kid, err := token.NewKID(rand)

		Convey("Then I should a valid KID", func() {
			So(err, ShouldBeNil)
			So(string(kid), ShouldNotBeBlank)
		})
	})

	Convey("When I generate a bad KID", t, func() {
		rand := rand.NewGenerator(&test.BadReader{})
		kid, err := token.NewKID(rand)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(string(kid), ShouldBeBlank)
		})
	})
}

func TestJWT(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())
	kid, _ := token.NewKID(rand)
	a, _ := ed25519.NewSigner(test.NewEd25519())
	jwt := token.NewJWT(kid, a)

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

			Convey("Then I should have a errror", func() {
				So(err, ShouldBeError)
			})
		})
	}
}
