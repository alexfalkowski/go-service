package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoder(t *testing.T) {
	for _, k := range []string{"yaml", "yml", "toml", "proto", "gob"} {
		Convey("When I get an encoder", t, func() {
			e := test.Encoder.Get(k)

			Convey("Then I should have an encoder", func() {
				So(e, ShouldNotBeNil)
			})
		})
	}

	for _, k := range []string{"test", "bob"} {
		Convey("When I get an encoder", t, func() {
			m := test.Encoder.Get(k)

			Convey("Then I should have no encoder", func() {
				So(m, ShouldBeNil)
			})
		})
	}
}
