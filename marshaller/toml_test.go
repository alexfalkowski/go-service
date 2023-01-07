package marshaller_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTOML(t *testing.T) {
	Convey("Given I have TOML marshaller", t, func() {
		m := marshaller.NewTOML()
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
}
