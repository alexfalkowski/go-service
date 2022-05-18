package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg := &test.NewTransportConfig().NSQ

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			producer := nsq.NewProducer(
				nsq.ProducerParams{Lifecycle: lc, Config: cfg, Marshaller: marshaller.NewMsgPack()},
				nsq.WithProducerLogger(logger), nsq.WithProducerTracer(tracer), nsq.WithProducerRetry(), nsq.WithProducerBreaker(),
				nsq.WithProducerMetrics(prometheus.NewProducerMetrics(lc, test.Version)),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(producer, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
