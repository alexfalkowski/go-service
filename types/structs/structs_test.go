package structs_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/types/structs"
	"github.com/alexfalkowski/go-service/types/validator"
	. "github.com/smartystreets/goconvey/convey"
)

type Config struct {
	Address string `validate:"hostname_port"`
}

func TestValid(t *testing.T) {
	Convey("Given I invalid struct", t, func() {
		structs.Register(validator.NewValidator())

		cfg := &Config{Address: "what?"}

		Convey("When I validate without context", func() {
			err := structs.Validate(cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I validated with context", func() {
			err := structs.ValidateWithContext(t.Context(), cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
