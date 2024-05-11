package marshaller_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestMap(t *testing.T) {
	for _, k := range []string{"yaml", "yml", "toml", "proto", "gob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create marshaller", func() {
				m := test.Marshaller.Get(k)

				Convey("Then I should have valid marshaller", func() {
					So(m, ShouldNotBeNil)
				})
			})
		})
	}

	for _, k := range []string{"test", "bob"} {
		Convey("Given I have map", t, func() {
			Convey("When I create marshaller", func() {
				m := test.Marshaller.Get(k)

				Convey("Then I should have none", func() {
					So(m, ShouldNotBeNil)
				})
			})
		})

		Convey("Given I have create a marshaller", t, func() {
			m := test.Marshaller.Get(k)

			Convey("When I marshal the data", func() {
				s := []byte("hello")

				d, err := m.Marshal(s)
				So(err, ShouldBeNil)

				Convey("Then I should be able to unmarshal", func() {
					err := m.Unmarshal(d, nil)
					So(err, ShouldBeNil)
				})
			})
		})
	}
}
