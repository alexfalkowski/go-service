package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/metrics/prometheus"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
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

			version := version.Version("1.0.0")

			tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
			So(err, ShouldBeNil)

			producer := nsq.NewProducer(
				nsq.ProducerParams{Lifecycle: lc, Config: cfg, Marshaller: marshaller.NewMsgPack(), Version: version},
				nsq.WithProducerLogger(logger), nsq.WithProducerTracer(tracer), nsq.WithProducerRetry(), nsq.WithProducerBreaker(),
				nsq.WithProducerMetrics(prometheus.NewProducerMetrics(lc, version)),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(producer, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
