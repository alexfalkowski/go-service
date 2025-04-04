package ed25519_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerator(t *testing.T) {
	Convey("Given I have an erroneous generator", t, func() {
		gen := ed25519.NewGenerator(rand.NewGenerator(rand.NewReader()))

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
		gen := ed25519.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))

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
	Convey("Given I have an signer", t, func() {
		cfg := test.NewEd25519()

		signer, err := ed25519.NewSigner(cfg)
		So(err, ShouldBeNil)

		verifier, err := ed25519.NewVerifier(cfg)
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := signer.Sign([]byte("test"))

			Convey("Then I should have verified the data", func() {
				So(verifier.Verify(e, []byte("test")), ShouldBeNil)
			})

			Convey("Then I should have keys", func() {
				So(signer.PrivateKey, ShouldNotBeNil)
				So(verifier.PublicKey, ShouldNotBeNil)
			})
		})
	})

	Convey("When I try to create a signer with no configuration", t, func() {
		signer, err := ed25519.NewSigner(nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no signer", func() {
			So(signer, ShouldBeNil)
		})
	})

	Convey("When I try to create a verifier with no configuration", t, func() {
		signer, err := ed25519.NewVerifier(nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no signer", func() {
			So(signer, ShouldBeNil)
		})
	})
}

//nolint:funlen
func TestInvalid(t *testing.T) {
	configs := []*ed25519.Config{
		{},
		{Public: test.Path("secrets/ed25519_public"), Private: test.Path("secrets/ed25519_private_invalid")},
	}

	for _, config := range configs {
		Convey("When I create a signer", t, func() {
			_, err := ed25519.NewSigner(config)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	configs = []*ed25519.Config{
		{},
		{Public: test.Path("secrets/ed25519_public_invalid"), Private: test.Path("secrets/ed25519_private")},
	}

	for _, config := range configs {
		Convey("When I create a signer", t, func() {
			_, err := ed25519.NewVerifier(config)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("Given I have an signer", t, func() {
		cfg := test.NewEd25519()

		signer, err := ed25519.NewSigner(cfg)
		So(err, ShouldBeNil)

		verifier, err := ed25519.NewVerifier(cfg)
		So(err, ShouldBeNil)

		Convey("When I sign the data", func() {
			sig, err := signer.Sign([]byte("test"))
			So(err, ShouldBeNil)

			sig = append(sig, byte('w'))

			Convey("Then I should have an error", func() {
				So(verifier.Verify(sig, []byte("test")), ShouldBeError)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		cfg := test.NewEd25519()

		signer, err := ed25519.NewSigner(cfg)
		So(err, ShouldBeNil)

		verifier, err := ed25519.NewVerifier(cfg)
		So(err, ShouldBeNil)

		Convey("When I sign one message", func() {
			e, err := signer.Sign([]byte("test"))
			So(err, ShouldBeNil)

			Convey("Then I comparing another message will gave an error", func() {
				So(verifier.Verify(e, []byte("bob")), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})

	Convey("When I create an verifier with an invalid public key", t, func() {
		_, err := ed25519.NewVerifier(&ed25519.Config{
			Public:  test.Path("secrets/rsa_public"),
			Private: test.Path("secrets/ed25519_private"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I create an signer with an invalid private key", t, func() {
		_, err := ed25519.NewSigner(&ed25519.Config{
			Public:  test.Path("secrets/ed25519_public"),
			Private: test.Path("secrets/rsa_private"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
