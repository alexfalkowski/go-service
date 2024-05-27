package hmac_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

//nolint:funlen
func TestAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		key, err := hmac.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(string(key), ShouldNotBeBlank)
		})
	})

	Convey("Given I have generated a key", t, func() {
		Convey("When I create an algo", func() {
			a, err := hmac.NewAlgo(test.NewHMAC())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Generate("test")

			Convey("Then I should compared the data", func() {
				So(a.Compare(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Generate("test")
			e += "wha"

			Convey("Then I should have an error", func() {
				So(a.Compare(e, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e := a.Generate("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(a.Compare(e, "bob"), ShouldBeError, errors.ErrMismatch)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		a, err := hmac.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Generate("test")

			Convey("Then I should compared the data", func() {
				So(a.Compare(e, "test"), ShouldBeNil)
			})
		})
	})
}
