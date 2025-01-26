package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	Convey("When I try to create a tracer with an invalid fs", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		_, err := tracer.NewTracer(lc, test.Environment, test.Version, test.Name, &test.ErrFS{}, test.NewOTLPTracerConfig(), logger)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
