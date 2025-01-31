package validator_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/validator"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

type Config struct {
	Address string `validate:"required,hostname_port"`
}

func TestValidator(t *testing.T) {
	Convey("Given I invalid struct", t, func() {
		cfg := &Config{Address: "what?"}

		Convey("When I validated it", func() {
			err := validator.ValidateStruct(cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I invalid field", t, func() {
		addr := "what?"

		Convey("When I validated it", func() {
			err := validator.ValidateField(&addr, "required,hostname_port")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
