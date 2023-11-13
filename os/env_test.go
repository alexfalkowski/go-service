package os_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/os"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEnv(t *testing.T) {
	Convey("When I retrieve env:HOME", t, func() {
		home := os.GetFromEnv("env:HOME")

		Convey("Then I should have a value for env:HOME", func() {
			So(home, ShouldNotBeEmpty)
			So(home, ShouldNotEqual, "env:HOME")
		})
	})

	Convey("When I retrieve env:BOB", t, func() {
		home := os.GetFromEnv("bob")

		Convey("Then I should have a value for env:HOME", func() {
			So(home, ShouldNotBeEmpty)
			So(home, ShouldEqual, "bob")
		})
	})
}
