package yaml_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	Convey("Given I have YAML encoder", t, func() {
		e := yaml.NewEncoder()

		b := test.Pool.Get()
		defer test.Pool.Put(b)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the YAML", func() {
			err := e.Encode(b, msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(b.String())

			Convey("Then I should have valid YAML", func() {
				So(s, ShouldEqual, "test: test")
			})
		})
	})

	Convey("Given I have YAML encoder", t, func() {
		e := yaml.NewEncoder()

		Convey("When I decode the YAML", func() {
			var msg map[string]string

			err := e.Decode(bytes.NewReader([]byte("test: test")), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
