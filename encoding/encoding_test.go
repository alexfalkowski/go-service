package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
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
			e := test.Encoder.Get(k)

			Convey("Then I should have an encoder", func() {
				So(e, ShouldNotBeNil)
			})
		})
	}
}
