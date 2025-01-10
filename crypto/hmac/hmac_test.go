package hmac_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		key, err := hmac.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(key, ShouldNotBeBlank)
		})
	})

	Convey("Given I have generated a key", t, func() {
		Convey("When I create an algo", func() {
			algo, err := hmac.NewAlgo(test.NewHMAC())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(algo, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e, err := algo.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I should compared the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		algo, err := hmac.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e, err := algo.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I should compared the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})
		})
	})
}

func TestInvalidAlgo(t *testing.T) {
	Convey("Given I have an algo", t, func() {
		algo, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			sign, err := algo.Sign("test")
			So(err, ShouldBeNil)

			sign += "wha"

			Convey("Then I should have an error", func() {
				So(algo.Verify(sign, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := hmac.NewAlgo(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e, err := algo.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I comparing another message will gave an error", func() {
				So(algo.Verify(e, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})
}
