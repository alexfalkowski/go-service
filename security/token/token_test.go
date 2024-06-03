package token_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/security/token"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerate(t *testing.T) {
	Convey("When I generate a key and hash", t, func() {
		k, h, err := token.Generate()
		So(err, ShouldBeNil)

		Convey("Then I should have a unauthenticated reply", func() {
			So(string(k), ShouldNotBeBlank)
			So(string(h), ShouldNotBeBlank)
		})
	})
}
