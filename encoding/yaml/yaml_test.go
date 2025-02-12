package yaml_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/yaml"
	"github.com/alexfalkowski/go-service/internal/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestEncoder(t *testing.T) {
	Convey("Given I have YAML encoder", t, func() {
		encoder := yaml.NewEncoder()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := map[string]string{"test": "test"}

		Convey("When I encode the YAML", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(bytes.String())

			Convey("Then I should have valid YAML", func() {
				So(s, ShouldEqual, "test: test")
			})
		})
	})

	Convey("Given I have YAML encoder", t, func() {
		encoder := yaml.NewEncoder()

		Convey("When I decode the YAML", func() {
			var msg map[string]string

			err := encoder.Decode(bytes.NewReader([]byte("test: test")), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
