package nsq_test

import (
	"testing"

	gn "github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/nsq"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg := &test.NewInsecureTransportConfig().NSQ

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)
			logger := test.NewLogger(lc)

			tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewDefaultTracerConfig(), Version: test.Version})
			So(err, ShouldBeNil)

			m, err := metrics.NewMeter(lc, test.Environment, test.Version)
			So(err, ShouldBeNil)

			producer, err := nsq.NewProducer(lc, cfg.Host, gn.NewMsgPackMarshaller(),
				nsq.WithProducerLogger(logger), nsq.WithProducerTracer(tracer), nsq.WithProducerRetry(&cfg.Retry),
				nsq.WithProducerBreaker(), nsq.WithProducerMetrics(m), nsq.WithProducerUserAgent(cfg.UserAgent),
			)
			So(err, ShouldBeNil)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(producer, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
