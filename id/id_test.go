package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidID(t *testing.T) {
	configs := []*id.Config{
		{Kind: "uuid"},
		{Kind: "ksuid"},
		{Kind: "nanoid"},
		{Kind: "ulid"},
		{Kind: "xid"},
	}

	for _, config := range configs {
		Convey("Given I have a generator", t, func() {
			gen, err := id.NewGenerator(config, test.Generators)
			So(err, ShouldBeNil)

			Convey("When I generate an id", func() {
				id := gen.Generate()

				Convey("Then I should an id", func() {
					So(id, ShouldNotBeBlank)
				})
			})
		})
	}
}

func TestInvalidID(t *testing.T) {
	Convey("When I create a generator with a nil config", t, func() {
		gen, err := id.NewGenerator(nil, test.Generators)
		So(err, ShouldBeNil)

		Convey("Then I should not have a generator", func() {
			So(gen, ShouldBeNil)
		})
	})

	Convey("When I create a generator with an invalid config", t, func() {
		_, err := id.NewGenerator(&id.Config{Kind: "invalid"}, test.Generators)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
