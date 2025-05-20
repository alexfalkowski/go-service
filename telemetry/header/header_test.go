package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSecrets(t *testing.T) {
	Convey("When I try to get secrets with an erroneous fs", t, func() {
		m := header.Map{"test": "none"}
		err := m.Secrets(test.ErrFS)

		Convey("Then I should have error", func() {
			So(err, ShouldBeError)
		})
	})
}
