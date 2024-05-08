package argon2_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/errors"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestHash(t *testing.T) {
	Convey("Given I have a hash", t, func() {
		h := argon2.NewAlgo()

		Convey("When I generate a hash", func() {
			s, err := h.Generate("test")
			So(err, ShouldBeNil)

			Convey("Then I should a hash", func() {
				So(s, ShouldNotBeBlank)
			})
		})

		Convey("When I generate a hash for test", func() {
			s, err := h.Generate("test")
			So(err, ShouldBeNil)

			Convey("Then I should a hash that is equal to test", func() {
				So(h.Compare(s, "test"), ShouldBeNil)
			})
		})

		Convey("When I generate a hash with the word steve", func() {
			s, err := h.Generate("steve")
			So(err, ShouldBeNil)

			Convey("Then comparing to bob should fail", func() {
				So(h.Compare(s, "bob"), ShouldBeError, errors.ErrMismatch)
			})
		})

		Convey("When I compare a non hashed value", func() {
			err := h.Compare("steve", "bob")

			Convey("Then comparing to bob should fail", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
