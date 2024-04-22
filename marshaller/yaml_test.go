package marshaller_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestYAML(t *testing.T) {
	Convey("Given I have YAML marshaller", t, func() {
		m := marshaller.NewYAML()
		msg := map[string]string{"test": "test"}

		Convey("When I marshall the YAML", func() {
			b, err := m.Marshal(msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(string(b))

			Convey("Then I should have valid YAML", func() {
				So(s, ShouldEqual, "test: test")
			})
		})
	})

	Convey("Given I have YAML marshaller", t, func() {
		m := marshaller.NewYAML()

		Convey("When I unmarshall the YAML", func() {
			var msg map[string]string

			err := m.Unmarshal([]byte("test: test"), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
