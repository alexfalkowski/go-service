package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	for _, k := range []string{"yaml", "yml", "toml", "proto", "gob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create encoder", func() {
				e := test.Encoder.Get(k)

				Convey("Then I should have valid encoder", func() {
					So(e, ShouldNotBeNil)
				})
			})
		})
	}

	for _, k := range []string{"test", "bob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create encoder", func() {
				m := test.Encoder.Get(k)

				Convey("Then I should have none", func() {
					So(m, ShouldBeNil)
				})
			})
		})
	}
}
