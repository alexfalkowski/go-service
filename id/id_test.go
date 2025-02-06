package id_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/id"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestID(t *testing.T) {
	configs := []*id.Config{
		nil,
		{},
		{Kind: "uuid"},
		{Kind: "ksuid"},
		{Kind: "nanoid"},
		{Kind: "ulid"},
		{Kind: "xid"},
	}

	for _, config := range configs {
		Convey("Given I have a generator", t, func() {
			gen := id.NewGenerator(config)

			Convey("When I generate an id", func() {
				id := gen.Generate()

				Convey("Then I should an id", func() {
					So(id, ShouldNotBeBlank)
				})
			})
		})
	}
}
