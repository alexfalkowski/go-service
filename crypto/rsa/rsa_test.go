package rsa_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		pub, pri, err := rsa.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(pub, ShouldNotBeBlank)
			So(pri, ShouldNotBeBlank)
		})
	})

	Convey("Given I have generated a key pair", t, func() {
		Convey("When I create an algo", func() {
			a, err := rsa.NewAlgo(test.NewRSA())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(a, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := rsa.NewAlgo(test.NewRSA())
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
		a, err := rsa.NewAlgo(nil)
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
	Convey("When I create an invalid algo", t, func() {
		a, err := rsa.NewAlgo(&rsa.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(a, ShouldBeNil)
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := rsa.NewAlgo(test.NewRSA())
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
		a, err := rsa.NewAlgo(test.NewRSA())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := a.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
