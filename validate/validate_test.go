package validate_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/validate"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

type Config struct {
	Address string `validate:"hostname_port"`
}

func TestValidator(t *testing.T) {
	Convey("Given I invalid struct", t, func() {
		cfg := &Config{Address: "what?"}

		Convey("When I validated it", func() {
			err := validate.Struct(cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I invalid field", t, func() {
		addr := "what?"

		Convey("When I validated it", func() {
			err := validate.Field(&addr, "hostname_port")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
