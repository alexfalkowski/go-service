package toml_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	t.Parallel()

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

			err := encoder.Decode(bytes.NewReader([]byte(`test = "test"`)), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
