package aes_test

import (
	"encoding/base64"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenertor(t *testing.T) {
	Convey("Given I have an erroneous generator", t, func() {
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
		cipher, err := aes.NewCipher(rand, test.NewAES())

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := cipher.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := cipher.Decrypt(enc)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})

	Convey("When I try to create a cipher with no config", t, func() {
		cipher, err := aes.NewCipher(nil, nil)
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

		cipher, err := aes.NewCipher(rand, &aes.Config{Key: test.Path("secrets/hooks")})
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			_, err := cipher.Encrypt("test")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I decrypt data", func() {
			m := base64.StdEncoding.EncodeToString([]byte("test"))
			_, err := cipher.Decrypt(m)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an cipher with an erroneous rand", t, func() {
		rand := rand.NewGenerator(&test.ErrReaderCloser{})

		cipher, err := aes.NewCipher(rand, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I try to encrypt data", func() {
			_, err := cipher.Encrypt("test")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	rand := rand.NewGenerator(rand.NewReader())

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := cipher.Encrypt("test")
			So(err, ShouldBeNil)

			enc += "wha"

			Convey("Then I should have an error", func() {
				_, err := cipher.Decrypt(enc)
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := aes.NewCipher(rand, test.NewAES())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := cipher.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
