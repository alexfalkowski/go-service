package hmac_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenertor(t *testing.T) {
	Convey("Given I have an erroneous generator", t, func() {
		gen := hmac.NewGenerator(rand.NewGenerator(rand.NewReader()))

		Convey("When I generate key", func() {
			key, err := gen.Generate()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(key, ShouldNotBeBlank)
			})
		})
	})

	Convey("Given I have an erroneous generator", t, func() {
		gen := hmac.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))

		Convey("When I generate key", func() {
			key, err := gen.Generate()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(key, ShouldBeBlank)
			})
		})
	})
}

func TestValidSigner(t *testing.T) {
	Convey("Given I have generated a key", t, func() {
		Convey("When I create an signer", func() {
			signer, err := hmac.NewSigner(test.NewHMAC())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(signer, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e, err := signer.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I should compared the data", func() {
				So(signer.Verify(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("When I create a signer with no configuration", t, func() {
		signer, err := hmac.NewSigner(nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no signer", func() {
			So(signer, ShouldBeNil)
		})
	})
}

func TestInvalidSigner(t *testing.T) {
	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			sign, err := signer.Sign("test")
			So(err, ShouldBeNil)

			sign += "wha"

			Convey("Then I should have an error", func() {
				So(signer.Verify(sign, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e, err := signer.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I comparing another message will gave an error", func() {
				So(signer.Verify(e, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})
}
