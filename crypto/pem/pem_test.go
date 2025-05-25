package pem_test

import (
	"errors"
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDecode(t *testing.T) {
	Convey("When I decode invalid path", t, func() {
		_, err := test.PEM.Decode(test.FilePath("none"), "n/a")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("When I decode invalid block", t, func() {
		_, err := test.PEM.Decode(test.FilePath("secrets/redis"), "n/a")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, pem.ErrInvalidBlock), ShouldBeTrue)
		})
	})

	Convey("When I decode invalid kind", t, func() {
		_, err := test.PEM.Decode(test.FilePath("secrets/rsa_public"), "what")

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, pem.ErrInvalidKind), ShouldBeTrue)
		})
	})
}
