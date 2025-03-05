package tracer_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestInvalidReader(t *testing.T) {
	Convey("When I try to create a tracer with an invalid fs", t, func() {
		lc := fxtest.NewLifecycle(t)
		params := tracer.Params{
			Lifecycle:   lc,
			Environment: test.Environment,
			Name:        test.Name,
			Version:     test.Version,
			FileSystem:  &test.ErrFS{},
			Config:      test.NewOTLPTracerConfig(),
		}
		_, err := tracer.NewTracer(params)

		Convey("Then I should have an error", func() {
			So(err, ShouldBeError)
		})
	})
}
