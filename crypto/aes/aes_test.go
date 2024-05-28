package aes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		key, err := aes.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(string(key), ShouldNotBeBlank)
		})
	})

	Convey("Given I have generated a key", t, func() {
		Convey("When I create an algo", func() {
			a, err := aes.NewAlgo(test.NewAES())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			e, err := a.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := a.Decrypt(e)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		a, err := aes.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			e, err := a.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := a.Decrypt(e)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})
}

func TestInvalidAlgo(t *testing.T) {
	Convey("Given I have an algo", t, func() {
		a, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			e, err := a.Encrypt("test")
			So(err, ShouldBeNil)

			e += "wha"

			Convey("Then I should have an error", func() {
				_, err := a.Decrypt(e)
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := a.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
