package nsq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gn "github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/test"
	tnsq "github.com/alexfalkowski/go-service/transport/nsq"
	ntracer "github.com/alexfalkowski/go-service/transport/nsq/telemetry/tracer"
	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func init() {
	tracer.Register()
}

//nolint:dupl
func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &test.NewTransportConfig().NSQ
		handler := test.NewConsumer(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(lc, "topic", "channel", cfg, handler, gn.NewMsgPackMarshaller(),
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
				tnsq.WithConsumerMetrics(m),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

//nolint:dupl
func TestInvalidConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &test.NewTransportConfig().NSQ
		handler := test.NewConsumer(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(lc, "schäfer", "schäfer", cfg, handler, gn.NewMsgPackMarshaller(),
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
				tnsq.WithConsumerMetrics(m),
			)

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidConsumerConfig(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &tnsq.Config{LookupHost: "invalid_host"}
		handler := test.NewConsumer(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(lc, "topic", "channel", cfg, handler, gn.NewMsgPackMarshaller(),
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
				tnsq.WithConsumerMetrics(m),
			)

			lc.RequireStart()

			Convey("Then I should have an error", func() {
				So(err, ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestReceiveMessage(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &test.NewTransportConfig().NSQ
		handler := test.NewConsumer(nil)

		err = tnsq.RegisterConsumer(lc, "topic", "channel", cfg, handler, gn.NewMsgPackMarshaller(),
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
			tnsq.WithConsumerMetrics(m),
		)
		So(err, ShouldBeNil)

		producer, err := tnsq.NewProducer(lc, cfg, gn.NewMsgPackMarshaller(),
			tnsq.WithProducerLogger(logger), tnsq.WithProducerTracer(tracer), tnsq.WithProducerRetry(), tnsq.WithProducerBreaker(),
			tnsq.WithProducerMetrics(m),
		)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := gn.NewMessage([]byte("test"))
			err = producer.Produce(context.Background(), "topic", message)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestReceiveMessageWithDefaultProducer(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &test.NewTransportConfig().NSQ
		handler := test.NewConsumer(nil)

		err = tnsq.RegisterConsumer(lc, "topic", "channel", cfg, handler, gn.NewMsgPackMarshaller(),
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
			tnsq.WithConsumerMetrics(m),
		)
		So(err, ShouldBeNil)

		producer, err := nsq.NewProducer(cfg.Host, nsq.NewConfig())
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			err = producer.Publish("topic", []byte("test"))
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should not receive a message", func() {
				So(handler.Message(), ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestReceiveError(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)

		tracer, err := ntracer.NewTracer(ntracer.Params{Lifecycle: lc, Config: test.NewTracerConfig(), Version: test.Version})
		So(err, ShouldBeNil)

		m, err := metrics.NewMeter(lc, test.Environment, test.Version)
		So(err, ShouldBeNil)

		cfg := &test.NewTransportConfig().NSQ
		handler := test.NewConsumer(errors.New("something went wrong"))

		err = tnsq.RegisterConsumer(lc, "topic", "channel", cfg, handler, gn.NewMsgPackMarshaller(),
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
			tnsq.WithConsumerMetrics(m),
		)
		So(err, ShouldBeNil)

		producer, err := tnsq.NewProducer(lc, cfg, gn.NewMsgPackMarshaller(),
			tnsq.WithProducerLogger(logger), tnsq.WithProducerTracer(tracer), tnsq.WithProducerRetry(), tnsq.WithProducerBreaker(),
			tnsq.WithProducerMetrics(m),
		)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := gn.NewMessage([]byte("test"))
			err = producer.Produce(context.Background(), "topic", message)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
