package compress_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestMap(t *testing.T) {
	for _, kind := range []string{"zstd", "s2", "snappy", "none"} {
		Convey("Given I have map", t, func() {
			Convey("When I create compressor", func() {
				cmp := test.Compressor.Get(kind)

				Convey("Then I should have valid compressor", func() {
					So(cmp, ShouldNotBeNil)
				})
			})
		})

		Convey("Given I have create a compressor", t, func() {
			cmp := test.Compressor.Get(kind)

			Convey("When I compress the data", func() {
				data := []byte("hello")
				d := cmp.Compress(data)

				Convey("Then I should have the same decompressed data", func() {
					ns, err := cmp.Decompress(d)
					So(err, ShouldBeNil)

					So(ns, ShouldEqual, data)
				})
			})
		})
	}

	for _, k := range []string{"test", "bob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create compressor", func() {
				cmp := test.Compressor.Get(k)

				Convey("Then I should have none", func() {
					So(cmp, ShouldNotBeNil)
				})
			})
		})

		Convey("Given I have create a compressor", t, func() {
			cmp := test.Compressor.Get(k)

			Convey("When I compress the data", func() {
				data := []byte("hello")
				d := cmp.Compress(data)

				Convey("Then I should have the same decompressed data", func() {
					ns, err := cmp.Decompress(d)
					So(err, ShouldBeNil)

					So(ns, ShouldEqual, data)
				})
			})
		})
	}
}
