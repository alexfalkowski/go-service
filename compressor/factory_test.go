package compressor_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidFactory(t *testing.T) {
	for _, k := range []string{"snappy"} {
		Convey("Given I have factory", t, func() {
			Convey("When I create compressor", func() {
				m, err := test.Compressor.Create(k)

				Convey("Then I should have valid marshaller", func() {
					So(err, ShouldBeNil)
					So(m, ShouldNotBeNil)
				})
			})
		})
	}
}

func TestInvalidFactory(t *testing.T) {
	for _, k := range []string{"test", "bob"} {
		Convey("Given I have factory", t, func() {
			Convey("When I create compressor", func() {
				m, err := test.Compressor.Create(k)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
					So(m, ShouldBeNil)
				})
			})
		})
	}
}
