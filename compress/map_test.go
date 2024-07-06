package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestMap(t *testing.T) {
	for _, k := range []string{"zstd", "s2", "snappy", "none"} {
		Convey("Given I have map", t, func() {
			Convey("When I create compressor", func() {
				m := test.Compressor.Get(k)

				Convey("Then I should have valid marshaller", func() {
					So(m, ShouldNotBeNil)
				})
			})
		})

		Convey("Given I have create a compressor", t, func() {
			m := test.Compressor.Get(k)

			Convey("When I compress the data", func() {
				s := []byte("hello")
				d := m.Compress(s)

				Convey("Then I should have the same decompressed data", func() {
					ns, err := m.Decompress(d)
					So(err, ShouldBeNil)

					So(ns, ShouldEqual, s)
				})
			})
		})
	}

	for _, k := range []string{"test", "bob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create compressor", func() {
				m := test.Compressor.Get(k)

				Convey("Then I should have none", func() {
					So(m, ShouldNotBeNil)
				})
			})
		})

		Convey("Given I have create a compressor", t, func() {
			m := test.Compressor.Get(k)

			Convey("When I compress the data", func() {
				s := []byte("hello")
				d := m.Compress(s)

				Convey("Then I should have the same decompressed data", func() {
					ns, err := m.Decompress(d)
					So(err, ShouldBeNil)

					So(ns, ShouldEqual, s)
				})
			})
		})
	}
}
