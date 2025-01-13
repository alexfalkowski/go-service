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
			So(pub, ShouldNotBeBlank)
			So(pri, ShouldNotBeBlank)
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := algo.Sign("test")

			Convey("Then I should have veirfied the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})

			Convey("Then I should have keys", func() {
				So(algo.PrivateKey(), ShouldNotBeNil)
				So(algo.PublicKey(), ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		algo, err := ed25519.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := algo.Sign("test")

			Convey("Then I should have veirfied the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})

			Convey("Then I should have missing keys", func() {
				So(algo.PrivateKey(), ShouldBeNil)
				So(algo.PublicKey(), ShouldBeNil)
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
		algo, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I sign the data", func() {
			e, _ := algo.Sign("test")
			e += "wha"

			Convey("Then I should have an error", func() {
				So(algo.Verify(e, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := ed25519.NewAlgo(test.NewEd25519())
		So(err, ShouldBeNil)

		Convey("When I sign one message", func() {
			e, _ := algo.Sign("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(algo.Verify(e, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})

	Convey("When I create an algo with an invalid public key", t, func() {
		_, err := ed25519.NewAlgo(&ed25519.Config{
			Public:  test.Path("secrets/rsa_public"),
			Private: test.Path("secrets/ed25519_private"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create an algo with an invalid private key", t, func() {
		_, err := ed25519.NewAlgo(&ed25519.Config{
			Public:  test.Path("secrets/ed25519_public"),
			Private: test.Path("secrets/rsa_private"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
