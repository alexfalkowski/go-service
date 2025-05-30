package rsa_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
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

func TestValid(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	Convey("When I create an encryptor", t, func() {
		cipher, err := rsa.NewEncryptor(rand, test.PEM, test.NewRSA())

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)
		})
	})

	Convey("When I create an decryptor", t, func() {
		cipher, err := rsa.NewDecryptor(rand, test.PEM, test.NewRSA())

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(cipher, ShouldNotBeNil)
		})
	})

	Convey("Given I have an cipher", t, func() {
		cfg := test.NewRSA()

		encryptor, err := rsa.NewEncryptor(rand, test.PEM, cfg)
		So(err, ShouldBeNil)

		decryptor, err := rsa.NewDecryptor(rand, test.PEM, cfg)
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			e, err := encryptor.Encrypt(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I should decrypt the data", func() {
				d, err := decryptor.Decrypt(e)
				So(err, ShouldBeNil)

				So(bytes.String(d), ShouldEqual, "test")
			})
		})
	})

	Convey("When I try to create a cipher with no configuration", t, func() {
		encryptor, err := rsa.NewEncryptor(rand, test.PEM, nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no encryptor", func() {
			So(encryptor, ShouldBeNil)
		})

		decryptor, err := rsa.NewDecryptor(rand, test.PEM, nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no decryptor", func() {
			So(decryptor, ShouldBeNil)
		})
	})
}

//nolint:funlen
func TestInvalid(t *testing.T) {
	rand := rand.NewGenerator(rand.NewReader())

	Convey("When I create an invalid encryptor", t, func() {
		encryptor, err := rsa.NewEncryptor(rand, test.PEM, &rsa.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(encryptor, ShouldBeNil)
		})
	})

	Convey("When I create an invalid decryptor", t, func() {
		decryptor, err := rsa.NewDecryptor(rand, test.PEM, &rsa.Config{})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(decryptor, ShouldBeNil)
		})
	})

	Convey("Given I have an cipher", t, func() {
		cfg := test.NewRSA()

		encryptor, err := rsa.NewEncryptor(rand, test.PEM, cfg)
		So(err, ShouldBeNil)

		decryptor, err := rsa.NewDecryptor(rand, test.PEM, cfg)
		So(err, ShouldBeNil)

		Convey("When I encrypt data", func() {
			enc, err := encryptor.Encrypt(strings.Bytes("test"))
			So(err, ShouldBeNil)

			enc = append(enc, byte('w'))

			Convey("Then I should have an error", func() {
				_, err := decryptor.Decrypt(enc)
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have an cipher", t, func() {
		cfg := test.NewRSA()

		decryptor, err := rsa.NewDecryptor(rand, test.PEM, cfg)
		So(err, ShouldBeNil)

		Convey("When I decrypt invalid data", func() {
			_, err := decryptor.Decrypt(strings.Bytes("test"))

			Convey("Then I have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("When I create an decryptor with an invalid public key", t, func() {
		config := &rsa.Config{
			Public:  test.FilePath("secrets/rsa_public"),
			Private: test.FilePath("secrets/ed25519_private"),
		}
		_, err := rsa.NewDecryptor(rand, test.PEM, config)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create an encryptor with an invalid private key", t, func() {
		config := &rsa.Config{
			Public:  test.FilePath("secrets/ed25519_public"),
			Private: test.FilePath("secrets/rsa_private"),
		}
		_, err := rsa.NewEncryptor(rand, test.PEM, config)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
