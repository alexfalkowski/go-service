package ssh_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidAlgo(t *testing.T) {
	Convey("When I generate keys", t, func() {
		pub, pri, err := ssh.Generate()

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeNil)
			So(pub, ShouldNotBeBlank)
			So(pri, ShouldNotBeBlank)
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := ssh.NewAlgo(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := algo.Sign("test")

			Convey("Then I should compared the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})
		})
	})

	Convey("Given I have a missing algo", t, func() {
		algo, err := ssh.NewAlgo(nil)
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := algo.Sign("test")

			Convey("Then I should compared the data", func() {
				So(algo.Verify(e, "test"), ShouldBeNil)
			})
		})
	})
}

func TestInvalidAlgo(t *testing.T) {
	Convey("When I create a algo", t, func() {
		_, err := ssh.NewAlgo(&ssh.Config{})

		Convey("Then I should not have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := ssh.NewAlgo(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign data", func() {
			e, _ := algo.Sign("test")
			e += "wha"

			Convey("Then I should have an error", func() {
				So(algo.Verify(e, "test"), ShouldBeError)
			})
		})
	})

	Convey("Given I have an algo", t, func() {
		algo, err := ssh.NewAlgo(test.NewSSH())
		So(err, ShouldBeNil)

		Convey("When I sign one message", func() {
			e, _ := algo.Sign("test")

			Convey("Then I comparing another message will gave an error", func() {
				So(algo.Verify(e, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})
	})
}
