package marshaller_test

import (
	"encoding/base64"
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestGob(t *testing.T) {
	Convey("Given I have gob marshaller", t, func() {
		m := marshaller.NewGOB()
		msg := map[string]string{"test": "test"}

		Convey("When I marshall the proto", func() {
			b, err := m.Marshal(msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(b)

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			})
		})
	})

	Convey("Given I have gob marshaller", t, func() {
		m := marshaller.NewGOB()

		Convey("When I unmarshall the gob", func() {
			b, err := base64.StdEncoding.DecodeString("DX8EAQL/gAABDAEMAAAO/4AAAQR0ZXN0BHRlc3Q=")
			So(err, ShouldBeNil)

			var msg map[string]string

			err = m.Unmarshal(b, &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg, ShouldEqual, map[string]string{"test": "test"})
			})
		})
	})
}
