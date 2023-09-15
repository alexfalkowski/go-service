package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	ntel "github.com/alexfalkowski/go-service/transport/nsq/telemetry"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	telemetry.RegisterTracer()
}

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg := &test.NewTransportConfig().NSQ

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			tracer, err := ntel.NewTracer(ntel.TracerParams{Lifecycle: lc, Config: test.NewTelemetryConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			producer := nsq.NewProducer(
				nsq.ProducerParams{Lifecycle: lc, Config: cfg, Marshaller: marshaller.NewMsgPack()},
				nsq.WithProducerLogger(logger), nsq.WithProducerTracer(tracer), nsq.WithProducerRetry(), nsq.WithProducerBreaker(),
				nsq.WithProducerMetrics(ntel.NewProducerMetrics(lc, test.Version)),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(producer, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
