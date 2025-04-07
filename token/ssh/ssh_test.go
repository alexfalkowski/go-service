package ssh_test

import (
	"testing"

	cs "github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	Convey("When I generate a SSH token", t, func() {
		token := ssh.NewToken(test.NewToken("ssh", "").SSH)

		tkn, err := token.Generate()
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			err := token.Verify(tkn)
			So(err, ShouldBeNil)
		})
	})

	Convey("When I try to create an ssh token", t, func() {
		ssh := ssh.NewToken(nil)

		Convey("Then I should not have a token", func() {
			So(ssh, ShouldBeNil)
		})
	})
}

func TestInvalid(t *testing.T) {
	Convey("When I generate a token with an invalid configuration", t, func() {
		cfg := &ssh.Config{
			Key: &ssh.Key{
				Name: "test",
				Config: &cs.Config{
					Public:  test.Path("secrets/ssh_public"),
					Private: test.Path("secrets/none"),
				},
			},
		}

		token := ssh.NewToken(cfg)
		_, err := token.Generate()

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	for _, tkn := range []string{"", "none-", "test-", "test-bob"} {
		token := ssh.NewToken(test.NewToken("ssh", "").SSH)

		Convey("When I verify a token", t, func() {
			err := token.Verify(tkn)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("When I verify a token with an invalid configuration", t, func() {
		cfg := &ssh.Config{
			Keys: ssh.Keys{
				&ssh.Key{
					Name: "test",
					Config: &cs.Config{
						Public:  test.Path("secrets/none"),
						Private: test.Path("secrets/ssh_private"),
					},
				},
			},
		}

		token := ssh.NewToken(cfg)
		err := token.Verify("test-bob")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
