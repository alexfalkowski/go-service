package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGenerator(t *testing.T) {
	Convey("Given I have an erroneous generator", t, func() {
		gen := ssh.NewGenerator(rand.NewGenerator(rand.NewReader()))

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
		gen := ssh.NewGenerator(rand.NewGenerator(&test.ErrReaderCloser{}))

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

func TestValidSigner(t *testing.T) {
	Convey("Given I have an signer", t, func() {
		signer, err := ssh.NewSigner(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := signer.Sign("test")

			Convey("Then I should compared the data", func() {
				So(signer.Verify(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("When I try to create a signer with no configuration", t, func() {
		signer, err := ssh.NewSigner(nil)
		So(err, ShouldBeNil)

		Convey("Then I should have no signer", func() {
			So(signer, ShouldBeNil)
		})
	})
}

//nolint:funlen
func TestInvalidSigner(t *testing.T) {
	Convey("When I create a signer", t, func() {
		_, err := ssh.NewSigner(&ssh.Config{})

		Convey("Then I should not have an ed25519 error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := ssh.NewSigner(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := signer.Sign("test")
			e += "what"

			Convey("Then I should have an error", func() {
				So(signer.Verify(e, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an signer", t, func() {
		signer, err := ssh.NewSigner(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign one message", func() {
			e, _ := signer.Sign("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(signer.Verify(e, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})

	Convey("When I have an invalid public key", t, func() {
		_, err := ssh.NewSigner(&ssh.Config{
			Public:  test.Path("secrets/redis"),
			Private: test.Path("secrets/ssh_private"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I have an invalid private key", t, func() {
		_, err := ssh.NewSigner(&ssh.Config{
			Public:  test.Path("secrets/ssh_public"),
			Private: test.Path("secrets/redis"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I have an missing private key", t, func() {
		_, err := ssh.NewSigner(&ssh.Config{
			Public:  test.Path("secrets/ssh_public"),
			Private: test.Path("secrets/none"),
		})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
