package pem_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/crypto/pem"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestDecode(t *testing.T) {
	Convey("When I decode invalid path", t, func() {
		_, err := pem.Decode("non existent", "n/a")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I decode invalid block", t, func() {
		_, err := pem.Decode(test.Path("secrets/redis"), "n/a")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, pem.ErrInvalidBlock), ShouldBeTrue)
		})
	})

	Convey("When I decode invalid kind", t, func() {
		_, err := pem.Decode(test.Path("secrets/rsa_public"), "what")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, pem.ErrInvalidKind), ShouldBeTrue)
		})
	})
}
