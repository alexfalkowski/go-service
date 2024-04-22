package marshaller_test

import (
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestJSON(t *testing.T) {
	Convey("Given I have JSON marshaller", t, func() {
		m := marshaller.NewJSON()
		msg := map[string]string{"test": "test"}

		Convey("When I marshall the JSON", func() {
			b, err := m.Marshal(msg)
			So(err, ShouldBeNil)

			s := strings.TrimSpace(string(b))

			Convey("Then I should have valid JSON", func() {
				So(s, ShouldEqual, "{\n    \"test\": \"test\"\n}")
			})
		})
	})

	Convey("Given I have JSON marshaller", t, func() {
		m := marshaller.NewJSON()

		Convey("When I unmarshall the JSON", func() {
			var msg map[string]string

			err := m.Unmarshal([]byte("{\n    \"test\": \"test\"\n}"), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid map", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
