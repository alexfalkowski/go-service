package rand_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestRand(t *testing.T) {
	Convey("When I generate random bytes", t, func() {
		c, err := rand.GenerateBytes(5)
		So(err, ShouldBeNil)

		Convey("Then I should have random bytes", func() {
			So(c, ShouldHaveLength, 5)
		})
	})

	Convey("When I generate random string", t, func() {
		s, err := rand.GenerateString(32)
		So(err, ShouldBeNil)

		Convey("Then I should have random bytes", func() {
			So(s, ShouldHaveLength, 32)
		})
	})
}
