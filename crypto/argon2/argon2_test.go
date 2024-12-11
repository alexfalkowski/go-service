package argon2_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/errors"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestAlgo(t *testing.T) {
	Convey("Given I have a hash", t, func() {
		h := argon2.NewAlgo()

		Convey("When I sign a hash", func() {
			s, err := h.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I should a hash", func() {
				So(s, ShouldNotBeBlank)
			})
		})

		Convey("When I sign a hash for test", func() {
			s, err := h.Sign("test")
			So(err, ShouldBeNil)

			Convey("Then I should a hash that is equal to test", func() {
				So(h.Verify(s, "test"), ShouldBeNil)
			})
		})

		Convey("When I sign a hash with the word steve", func() {
			s, err := h.Sign("steve")
			So(err, ShouldBeNil)

			Convey("Then verifying to bob should fail", func() {
				So(h.Verify(s, "bob"), ShouldBeError, errors.ErrInvalidMatch)
			})
		})

		Convey("When I compare a non hashed value", func() {
			err := h.Verify("steve", "bob")

			Convey("Then comparing to bob should fail", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
