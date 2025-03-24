package fields_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/types/fields"
	"github.com/alexfalkowski/go-service/types/validator"
	. "github.com/smartystreets/goconvey/convey"
)

type Config struct {
	Address string `validate:"hostname_port"`
}

func TestValid(t *testing.T) {
	Convey("Given I invalid field", t, func() {
		fields.Register(validator.NewValidator())

		addr := "what?"

		Convey("When I validated without context", func() {
			err := fields.Validate(&addr, "hostname_port")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I validate with context", func() {
			err := fields.ValidateWithContext(t.Context(), &addr, "hostname_port")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
