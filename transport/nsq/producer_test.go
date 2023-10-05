package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/telemetry/metrics/prometheus"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg := &test.NewTransportConfig().NSQ

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)
			tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			producer := nsq.NewProducer(lc, cfg, marshaller.NewMsgPack(),
				nsq.WithProducerLogger(logger), nsq.WithProducerTracer(tracer), nsq.WithProducerRetry(), nsq.WithProducerBreaker(),
				nsq.WithProducerMetrics(prometheus.NewProducerCollector(lc, test.Version)),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(producer, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
