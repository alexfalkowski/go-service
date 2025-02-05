package rsa_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenertor(t *testing.T) {
	Convey("Given I have an erroneous generator", t, func() {
		gen := rsa.NewGenerator(rand.NewGenerator(rand.NewReader()))

		Convey("When I generate keys", func() {
			pub, pri, err := gen.Generate()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(pub, ShouldNotBeBlank)
				So(pri, ShouldNotBeBlank)
			})
		})
	})

	Convey("Given I have an erroneous generator", t, func() {
		gen := rsa.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))

		Convey("When I generate keys", func() {
			pub, pri, err := gen.Generate()

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
				So(pub, ShouldBeBlank)
				So(pri, ShouldBeBlank)
			})
		})
	})
}

func TestValidCipher(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	Convey("Given I have generated a key pair", t, func() {
		Convey("When I create an cipher", func() {
			cipher, err := rsa.NewCipher(rand, test.NewRSA())

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
				So(cipher, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := rsa.NewCipher(rand, test.NewRSA())
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			e, err := cipher.Encrypt("test")
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := cipher.Decrypt(e)
				So(err, ShouldBeNil)

				So(d, ShouldEqual, "test")
			})
		})
	})

	Convey("When I try to create a cipher with no configuration", t, func() {
		cipher, err := rsa.NewCipher(rand, nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no cipher", func() {
			So(cipher, ShouldBeNil)
		})
	})
}

//nolint:funlen
func TestInvalidCipher(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	Convey("When I create an invalid cipher", t, func() {
		cipher, err := rsa.NewCipher(rand, &rsa.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(cipher, ShouldBeNil)
		})
	})

	Convey("Given I have an cipher", t, func() {
		cipher, err := rsa.NewCipher(rand, test.NewRSA())
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
		cipher, err := rsa.NewCipher(rand, test.NewRSA())
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := cipher.Decrypt("test")

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("When I create an cipher with an invalid public key", t, func() {
		config := &rsa.Config{
			Public:  test.Path("secrets/ed25519_public"),
			Private: test.Path("secrets/rsa_private"),
		}
		_, err := rsa.NewCipher(rand, config)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create an cipher with an invalid private key", t, func() {
		config := &rsa.Config{
			Public:  test.Path("secrets/rsa_public"),
			Private: test.Path("secrets/ed25519_private"),
		}
		_, err := rsa.NewCipher(rand, config)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
