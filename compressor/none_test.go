package compressor_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/compressor"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestNone(t *testing.T) {
	Convey("Given I have a compressor", t, func() {
		m := compressor.NewNone()
		msg := []byte("test")

		Convey("When I compress", func() {
			b := m.Compress(msg)

			Convey("Then I should the same value", func() {
				So(b, ShouldEqual, msg)
			})
		})
	})

	Convey("Given I have a compressor", t, func() {
		m := compressor.NewNone()
		msg := []byte("test")

		Convey("When I decompress", func() {
			b, err := m.Decompress(msg)
			So(err, ShouldBeNil)

			Convey("Then I should the same value", func() {
				So(b, ShouldEqual, msg)
			})
		})
	})
}
