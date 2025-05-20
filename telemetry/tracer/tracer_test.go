package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	Convey("When I try to create a tracer with an invalid fs", t, func() {
		lc := fxtest.NewLifecycle(t)
		params := tracer.Params{
			Lifecycle:   lc,
			ID:          test.ID,
			Name:        test.Name,
			Version:     test.Version,
			Environment: test.Environment,
			FileSystem:  test.ErrFS,
			Config:      test.NewOTLPTracerConfig(),
		}
		_, err := tracer.NewTracer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
