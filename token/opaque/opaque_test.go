package opaque_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/internal/test"
	te "github.com/alexfalkowski/go-service/token/errors"
	"github.com/alexfalkowski/go-service/token/opaque"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	cfg := test.NewToken("opaque")
	token := opaque.NewToken(cfg.Opaque, rand.NewGenerator(rand.NewReader()))

	Convey("When I generate an opaque token", t, func() {
		tkn := token.Generate()

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			err := token.Verify(tkn)
			So(err, ShouldBeNil)
		})

		Convey("Then I should have an error due to invalid token", func() {
			err := token.Verify("bob")
			So(err, ShouldBeError)
		})
	})
}

func TestInvalid(t *testing.T) {
	cfg := test.NewToken("opaque")
	token := opaque.NewToken(cfg.Opaque, rand.NewGenerator(rand.NewReader()))

	name := test.Name.String()
	keys := []string{
		"",
		"none_test_test",
		name + "_test_test",
		name + "_test_1",
	}

	for _, key := range keys {
		Convey("When I verify a token", t, func() {
			err := token.Verify(key)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, te.ErrInvalidMatch), ShouldBeTrue)
			})
		})
	}
}
