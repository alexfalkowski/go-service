package hooks_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/hooks"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestHooks(t *testing.T) {
	Convey("When I generate a secret", t, func() {
		c, err := hooks.Generate()
		So(err, ShouldBeNil)

		Convey("Then I should have random secret", func() {
			So(string(c), ShouldNotBeBlank)
		})
	})
}
