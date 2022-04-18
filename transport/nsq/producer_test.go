package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg := test.NewNSQConfig()

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)

			logger, err := zap.NewLogger(lc, zap.NewConfig())
			So(err, ShouldBeNil)

			tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
			So(err, ShouldBeNil)

			params := &nsq.ProducerParams{
				Lifecycle: lc,
				Config:    cfg,
				Logger:    logger,
				Tracer:    tracer,
			}
			_, err = nsq.NewProducer(params)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
