package marshaller_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey"
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
}
