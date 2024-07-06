package toml_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/toml"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestTOML(t *testing.T) {
	Convey("Given I have TOML marshaller", t, func() {
		m := toml.NewMarshaller()
		msg := map[string]string{"test": "test"}

		Convey("When I marshall the TOML", func() {
			b, err := m.Marshal(msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(string(b))

			Convey("Then I should have valid TOML", func() {
				So(s, ShouldEqual, `test = "test"`)
			})
		})
	})

	Convey("Given I have TOML marshaller", t, func() {
		m := toml.NewMarshaller()

		Convey("When I unmarshal the TOML", func() {
			var msg map[string]string

			err := m.Unmarshal([]byte(`test = "test"`), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
