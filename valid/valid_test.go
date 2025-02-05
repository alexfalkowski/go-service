package valid_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/valid"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

type Config struct {
	Address string `validate:"hostname_port"`
}

func TestValid(t *testing.T) {
	t.Parallel()

	Convey("Given I invalid struct", t, func() {
		cfg := &Config{Address: "what?"}

		Convey("When I validated it", func() {
			err := valid.Struct(cfg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I invalid field", t, func() {
		addr := "what?"

		Convey("When I validated it", func() {
			err := valid.Field(&addr, "hostname_port")

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
