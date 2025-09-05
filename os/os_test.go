package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/os"
	. "github.com/smartystreets/goconvey/convey"
)

//nolint:usetesting
func TestEnv(t *testing.T) {
	Convey("Given I have a env variable", t, func() {
		key := "__ENV_KEY"

		err := os.Setenv(key, "test")
		So(err, ShouldBeNil)

		Convey("When I get the value", func() {
			value := os.Getenv(key)

			Convey("Then I should have a value", func() {
				So(value, ShouldEqual, "test")
			})
		})

		err = os.Unsetenv(key)
		So(err, ShouldBeNil)
	})
}
