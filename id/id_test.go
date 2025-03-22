package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/id"
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
			gen, err := id.NewGenerator(config)
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
	configs := []*id.Config{
		nil,
		{},
	}

	for _, config := range configs {
		Convey("When I create a generator", t, func() {
			gen, err := id.NewGenerator(config)
			So(err, ShouldBeNil)

			Convey("Then I should not have a generator", func() {
				So(gen, ShouldBeNil)
			})
		})
	}

	Convey("When I create a generator", t, func() {
		_, err := id.NewGenerator(&id.Config{Kind: "invalid"})

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
