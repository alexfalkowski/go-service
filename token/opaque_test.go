package token_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestOpaque(t *testing.T) {
	opaque := token.NewOpaque(test.Name, rand.NewGenerator(rand.NewReader()))

	Convey("When I generate a JWT token", t, func() {
		token := opaque.Generate()

		Convey("Then I should have a token", func() {
			So(token, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			err := opaque.Verify(token)
			So(err, ShouldBeNil)
		})

		Convey("Then I should have an error due to invalid token", func() {
			err := opaque.Verify("bob")
			So(err, ShouldBeError)
		})
	})

	keys := []string{
		"",
		"none_test_test",
		string(test.Name) + "_test_test",
		string(test.Name) + "_test_1",
	}

	for _, key := range keys {
		Convey("When I verify a token", t, func() {
			err := opaque.Verify(key)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, token.ErrInvalidMatch), ShouldBeTrue)
			})
		})
	}
}
