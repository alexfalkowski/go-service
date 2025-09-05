package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValid(t *testing.T) {
	Convey("When I generate a SSH token", t, func() {
		token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

		tkn, err := token.Generate()
		So(err, ShouldBeNil)

		Convey("Then I should have a token", func() {
			So(tkn, ShouldNotBeBlank)
		})

		Convey("Then I should be able to verify the token", func() {
			sub, err := token.Verify(tkn)
			So(err, ShouldBeNil)
			So(sub, ShouldEqual, test.UserID.String())
		})
	})

	Convey("When I try to create an ssh token", t, func() {
		ssh := ssh.NewToken(nil, nil)

		Convey("Then I should not have a token", func() {
			So(ssh, ShouldBeNil)
		})
	})
}

func TestInvalid(t *testing.T) {
	Convey("When I generate a token with an invalid configuration", t, func() {
		cfg := &ssh.Config{
			Key: &ssh.Key{
				Name:   "test",
				Config: test.NewSSH("secrets/ssh_public", "secrets/none"),
			},
		}

		token := ssh.NewToken(cfg, test.FS)
		_, err := token.Generate()

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	for _, tkn := range []string{"", "none-", "test-", "test-bob"} {
		token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

		Convey("When I verify a token", t, func() {
			_, err := token.Verify(tkn)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	}

	Convey("When I verify a token with an invalid configuration", t, func() {
		cfg := &ssh.Config{
			Keys: ssh.Keys{
				&ssh.Key{
					Name:   "test",
					Config: test.NewSSH("secrets/none", "secrets/ssh_private"),
				},
			},
		}

		token := ssh.NewToken(cfg, test.FS)
		_, err := token.Verify("test-bob")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I try to create a ssh token", t, func() {
		token := ssh.NewToken(nil, test.FS)

		Convey("Then I should have no token", func() {
			So(token, ShouldBeNil)
		})
	})
}
