package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/security/header"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidParseAuthorization(t *testing.T) {
	Convey("Given I have a valid header", t, func() {
		h := "Bearer token"

		Convey("When I parse authorization", func() {
			t, c, err := header.ParseAuthorization(h)
			So(err, ShouldBeNil)

			Convey("Then I should have a valid type and credentials", func() {
				So(t, ShouldEqual, "Bearer")
				So(c, ShouldEqual, "token")
			})
		})
	})
}

func TestMissingParseAuthorization(t *testing.T) {
	Convey("Given I have a missing header", t, func() {
		h := ""

		Convey("When I parse authorization", func() {
			_, _, err := header.ParseAuthorization(h)

			Convey("Then I should have a invalid type and credentials", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, header.ErrInvalidAuthorization)
			})
		})
	})
}

func TestNotSupportedParseAuthorization(t *testing.T) {
	Convey("Given I have a not supported header", t, func() {
		h := "Bob token"

		Convey("When I parse authorization", func() {
			_, _, err := header.ParseAuthorization(h)

			Convey("Then I should have a not supported type and credentials", func() {
				So(err, ShouldBeError)
				So(err, ShouldEqual, header.ErrNotSupportedAuthorization)
			})
		})
	})
}
