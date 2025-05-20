package bcrypt_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v2/strings"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSigner(t *testing.T) {
	Convey("Given I have a hash", t, func() {
		signer := bcrypt.NewSigner()

		Convey("When I sign a hash", func() {
			s, err := signer.Sign(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I should a hash", func() {
				So(s, ShouldNotBeEmpty)
			})
		})

		Convey("When I sign a hash for test", func() {
			s, err := signer.Sign(strings.Bytes("test"))
			So(err, ShouldBeNil)

			Convey("Then I should a hash that is equal to test", func() {
				So(signer.Verify(s, strings.Bytes("test")), ShouldBeNil)
			})
		})

		Convey("When I sign a hash with the word steve", func() {
			s, err := signer.Sign(strings.Bytes("steve"))
			So(err, ShouldBeNil)

			Convey("Then verifying to bob should fail", func() {
				So(signer.Verify(s, strings.Bytes("bob")), ShouldBeError)
			})
		})

		Convey("When I compare a non hashed value", func() {
			err := signer.Verify(strings.Bytes("steve"), strings.Bytes("bob"))

			Convey("Then comparing to bob should fail", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
