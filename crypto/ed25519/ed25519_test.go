package ed25519_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		pub, pri, err := ed25519.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(string(pub), ShouldNotBeBlank)
			So(string(pri), ShouldNotBeBlank)
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Sign("test")

			Convey("Then I should compared the data", func() {
				So(a.Verify(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		a, err := ed25519.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Sign("test")

			Convey("Then I should compared the data", func() {
				So(a.Verify(e, "test"), ShouldBeNil)
			})
		})
	})
}

func TestInvalidAlgo(t *testing.T) {
	Convey("When I create a algo", t, func() {
		_, err := ed25519.NewAlgo(&ed25519.Config{})

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Sign("test")
			e += "wha"

			Convey("Then I should have an error", func() {
				So(a.Verify(e, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		a, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e := a.Sign("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(a.Verify(e, "bob"), ShouldBeError, errors.ErrMismatch)
			})
		})
	})
}
