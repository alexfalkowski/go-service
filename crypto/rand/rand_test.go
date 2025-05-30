package rand_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidRand(t *testing.T) {
	gen := rand.NewGenerator(rand.NewReader())

	Convey("When I generate random bytes", t, func() {
		c, err := gen.GenerateBytes(5)
		So(err, ShouldBeNil)

		Convey("Then I should have random bytes", func() {
			So(c, ShouldHaveLength, 5)
		})
	})

	Convey("When I generate random letters", t, func() {
		s, err := gen.GenerateText(32)
		So(err, ShouldBeNil)

		Convey("Then I should have random letters", func() {
			So(s, ShouldHaveLength, 32)
		})
	})
}

func TestInvalidRand(t *testing.T) {
	gen := rand.NewGenerator(&test.ErrReaderCloser{})

	Convey("When I generate random string with an erroneous reader", t, func() {
		_, err := gen.GenerateText(5)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
