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
			So(key, ShouldNotBeBlank)
		})
	})

	Convey("Given I have generated a key", t, func() {
		Convey("When I create an algo", func() {
			algo, err := aes.NewAlgo(test.NewAES())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(algo, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := algo.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := algo.Decrypt(enc)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		algo, err := aes.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := algo.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := algo.Decrypt(enc)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})
}

func TestInvalidAlgo(t *testing.T) {
	Convey("Given I have an algo", t, func() {
		algo, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := algo.Encrypt("test")
			So(err, ShouldBeNil)

			enc += "wha"

			Convey("Then I should have an error", func() {
				_, err := algo.Decrypt(enc)
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := aes.NewAlgo(test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := algo.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
