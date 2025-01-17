package header_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/telemetry/header"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestSecrets(t *testing.T) {
	Convey("When I try to get secrets with a bad fs", t, func() {
		m := header.Map{"test": "bad_path"}
		err := m.Secrets(&test.BadFS{})

		Convey("Then I should have error", func() {
			So(err, ShouldBeError)
		})
	})
}
