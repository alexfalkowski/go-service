package aes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/aes"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

//nolint:funlen
func TestAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		key, err := aes.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(string(key), ShouldNotBeBlank)
		})
	})

	Convey("When I create an invalid algo", t, func() {
		a, err := aes.NewAlgo(&aes.Config{Key: "==sdd"})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(a, ShouldNotBeNil)
		})
	})

	Convey("Given I have generated a key", t, func() {
		key, err := aes.Generate()
		So(err, ShouldBeNil)

		Convey("When I create an algo", func() {
			a, err := aes.NewAlgo(&aes.Config{Key: key})

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		key, err := aes.Generate()
		So(err, ShouldBeNil)

		a, err := aes.NewAlgo(&aes.Config{Key: key})
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

	Convey("Given I have an algo", t, func() {
		key, err := aes.Generate()
		So(err, ShouldBeNil)

		a, err := aes.NewAlgo(&aes.Config{Key: key})
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := a.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
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
