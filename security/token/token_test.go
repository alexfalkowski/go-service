package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerateConfig(t *testing.T) {
	Convey("When I generate a key and hash", t, func() {
		tkn := token.NewTokenizer(test.NewToken(), argon2.NewAlgo())

		k, h, err := tkn.GenerateConfig()
		So(err, ShouldBeNil)

		Convey("Then I should have a unauthenticated reply", func() {
			So(k, ShouldNotBeBlank)
			So(h, ShouldNotBeBlank)
		})
	})
}
