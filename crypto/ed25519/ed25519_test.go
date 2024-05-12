package ed25519_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/errors"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

//nolint:funlen
func TestAlgo(t *testing.T) {
	Convey("When I generate", t, func() {
		pub, pri, err := ed25519.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(string(pub), ShouldNotBeBlank)
			So(string(pri), ShouldNotBeBlank)
		})
	})

	Convey("When I create a algo", t, func() {
		a, err := ed25519.NewAlgo(&ed25519.Config{})

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(a.PrivateKey(), ShouldNotBeNil)
			So(a.PublicKey(), ShouldNotBeNil)
		})
	})

	Convey("When I create a missing algo", t, func() {
		a, err := ed25519.NewAlgo(nil)

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(a.PrivateKey(), ShouldNotBeNil)
			So(a.PublicKey(), ShouldNotBeNil)
		})
	})

	type tuple [2]string

	tus := []tuple{{"dGVzdAo==", "test"}, {"test", "dGVzdAo=="}}

	for _, tu := range tus {
		Convey("When I try to create a bad algo", t, func() {
			_, err := ed25519.NewAlgo(&ed25519.Config{Public: ed25519.PublicKey(tu[0]), Private: ed25519.PrivateKey(tu[1])})

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have an algo", t, func() {
		pub, pri, err := ed25519.Generate()
		So(err, ShouldBeNil)

		a, err := ed25519.NewAlgo(&ed25519.Config{Public: pub, Private: pri})
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Generate("test")

			Convey("Then I should compared the data", func() {
				So(a.Compare(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		pub, pri, err := ed25519.Generate()
		So(err, ShouldBeNil)

		a, err := ed25519.NewAlgo(&ed25519.Config{Public: pub, Private: pri})
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e := a.Generate("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(a.Compare(e, "bob"), ShouldBeError, errors.ErrMismatch)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		a, err := ed25519.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e := a.Generate("test")

			Convey("Then I should compared the data", func() {
				So(a.Compare(e, "test"), ShouldBeNil)
			})
		})
	})
}
