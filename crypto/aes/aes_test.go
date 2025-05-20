package aes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	Convey("Given I have a generator", t, func() {
		gen := aes.NewGenerator(rand.NewGenerator(rand.NewReader()))

		Convey("When I generate key", func() {
			key, err := gen.Generate()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(key, ShouldNotBeBlank)
			})
		})
	})

	Convey("Given I have a erroneous generator", t, func() {
		gen := aes.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))

		Convey("When I generate key", func() {
			key, err := gen.Generate()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(key, ShouldBeBlank)
			})
		})
	})
}

func TestValidCipher(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	Convey("When I create an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := cipher.Encrypt(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := cipher.Decrypt(enc)
				So(err, ShouldBeNil)

				So(bytes.String(d), ShouldEqual, "test")
			})
		})
	})

	Convey("When I try to create a cipher with no config", t, func() {
		cipher, err := aes.NewCipher(nil, nil, nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no cipher", func() {
			So(cipher, ShouldBeNil)
		})
	})
}

//nolint:funlen
func TestInvalidCipher(t *testing.T) {
	Convey("Given I have an cipher with invalid key", t, func() {
		rand := rand.NewGenerator(rand.NewReader())

		cipher, err := aes.NewCipher(rand, test.FS, &aes.Config{Key: test.Path("secrets/aes_invalid")})
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			_, err := cipher.Encrypt(strings.Bytes("test"))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I decrypt data", func() {
			_, err := cipher.Decrypt(strings.Bytes("test"))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an cipher with an erroneous rand", t, func() {
		rand := rand.NewGenerator(&test.ErrReaderCloser{})

		cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I try to encrypt data", func() {
			_, err := cipher.Encrypt(strings.Bytes("test"))

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	rand := rand.NewGenerator(rand.NewReader())

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := cipher.Encrypt(strings.Bytes("test"))
			So(err, ShouldBeNil)

			enc = append(enc, byte('w'))

			Convey("Then I should have an error", func() {
				_, err := cipher.Decrypt(enc)
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.FS, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := cipher.Decrypt(strings.Bytes("test"))

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
