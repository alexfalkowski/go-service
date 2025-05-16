package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHooks(t *testing.T) {
	gen := hooks.NewGenerator(rand.NewGenerator(rand.NewReader()))

	Convey("When I generate a secret", t, func() {
		c, err := gen.Generate()
		So(err, ShouldBeNil)

		Convey("Then I should have random secret", func() {
			So(c, ShouldNotBeBlank)
		})
	})

	Convey("When I create a hook with an missing secret", t, func() {
		_, err := hooks.New(test.FS, &hooks.Config{Secret: test.Path("secrets/none")})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create a hook with an invalid secret", t, func() {
		_, err := hooks.New(test.FS, &hooks.Config{Secret: test.Path("secrets/redis")})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create a hook with a missing config", t, func() {
		h, err := hooks.New(nil, nil)

		Convey("Then I should not have a hook", func() {
			So(h, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})
}
