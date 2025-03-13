package argon2_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSigner(t *testing.T) {
	Convey("Given I have a hash", t, func() {
		signer := argon2.NewSigner()

		Convey("When I sign a hash", func() {
			s, err := signer.Sign([]byte("test"))
			So(err, ShouldBeNil)

			Convey("Then I should a hash", func() {
				So(s, ShouldNotBeEmpty)
			})
		})

		Convey("When I sign a hash for test", func() {
			s, err := signer.Sign([]byte("test"))
			So(err, ShouldBeNil)

			Convey("Then I should a hash that is equal to test", func() {
				So(signer.Verify(s, []byte("test")), ShouldBeNil)
			})
		})

		Convey("When I sign a hash with the word steve", func() {
			s, err := signer.Sign([]byte("steve"))
			So(err, ShouldBeNil)

			Convey("Then verifying to bob should fail", func() {
				So(signer.Verify(s, []byte("bob")), ShouldBeError, errors.ErrInvalidMatch)
			})
		})

		Convey("When I compare a non hashed value", func() {
			err := signer.Verify([]byte("steve"), []byte("bob"))

			Convey("Then comparing to bob should fail", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
