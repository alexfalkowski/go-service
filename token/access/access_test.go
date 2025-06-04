package access_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/access"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewController(t *testing.T) {
	Convey("When I try to create an access controller with an invalid config", t, func() {
		_, err := access.NewController(&access.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I try to create an access controller with an missing config", t, func() {
		controller, err := access.NewController(nil)

		Convey("Then I should not have an error", func() {
			So(controller, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestHasAccess(t *testing.T) {
	Convey("Given I have an access controller", t, func() {
		config := test.NewAccessConfig()

		controller, err := access.NewController(config)
		So(err, ShouldBeNil)

		Convey("Which I check for access", func() {
			ok, err := controller.HasAccess("alice", "service", "read")
			So(err, ShouldBeNil)

			Convey("Then I should have access", func() {
				So(ok, ShouldBeTrue)
			})
		})
	})
}
