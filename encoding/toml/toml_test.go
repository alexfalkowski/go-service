package toml_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoder(t *testing.T) {
	Convey("Given I have TOML encoder", t, func() {
		encoder := toml.NewEncoder()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the TOML", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(bytes.String())

			Convey("Then I should have valid TOML", func() {
				So(s, ShouldEqual, `test = "test"`)
			})
		})
	})

	Convey("Given I have TOML encoder", t, func() {
		encoder := toml.NewEncoder()

		Convey("When I decode the TOML", func() {
			var msg map[string]string

			err := encoder.Decode(bytes.NewBufferString(`test = "test"`), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
