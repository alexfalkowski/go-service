package hmac_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	Convey("Given I have a generator", t, func() {
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
			signer, err := hmac.NewSigner(test.FS, test.NewHMAC())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(signer, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			e, err := signer.Sign(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I should compared the data", func() {
				So(signer.Verify(e, strings.Bytes("test")), ShouldBeNil)
			})
		})
	})

	Convey("When I create a signer with no configuration", t, func() {
		signer, err := hmac.NewSigner(nil, nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no signer", func() {
			So(signer, ShouldBeNil)
		})
	})
}

func TestInvalidSigner(t *testing.T) {
	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate data", func() {
			sign, err := signer.Sign(strings.Bytes("test"))
			So(err, ShouldBeNil)

			sign = append(sign, byte('w'))

			Convey("Then I should have an error", func() {
				So(signer.Verify(sign, strings.Bytes("test")), ShouldBeError)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := hmac.NewSigner(test.FS, test.NewHMAC())
		So(err, ShouldBeNil)

		Convey("When I generate one message", func() {
			e, err := signer.Sign(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I comparing another message will gave an error", func() {
				So(signer.Verify(e, strings.Bytes("bob")), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})
}
