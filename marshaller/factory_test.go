package marshaller_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/marshaller"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestValidFactory(t *testing.T) {
	for _, k := range []string{"yaml", "yml", "toml"} {
		Convey("Given I have factory", t, func() {
			p := marshaller.FactoryParams{
				JSON: marshaller.NewJSON(),
				TOML: marshaller.NewTOML(),
				YAML: marshaller.NewYAML(),
			}
			f := marshaller.NewFactory(p)

			Convey("When I create marshaller", func() {
				m, err := f.Create(k)

				Convey("Then I should have valid marshaller", func() {
					So(err, ShouldBeNil)
					So(m, ShouldNotBeNil)
				})
			})
		})
	}
}

func TestInvalidFactory(t *testing.T) {
	for _, k := range []string{"test", "bob"} {
		Convey("Given I have factory", t, func() {
			f := marshaller.NewFactory(marshaller.FactoryParams{})

			Convey("When I create marshaller", func() {
				m, err := f.Create(k)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
					So(m, ShouldBeNil)
				})
			})
		})
	}
}
