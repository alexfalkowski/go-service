package nsq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	tnsq "github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

// nolint:dupl
func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(
				tnsq.ConsumerParams{
					Lifecycle: lc, Topic: "topic", Channel: "channel", Config: cfg,
					Handler: handler, Marshaller: marshaller.NewMsgPack(),
					Version: version,
				},
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
			)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

// nolint:dupl
func TestInvalidConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewDatadogConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(
				tnsq.ConsumerParams{
					Lifecycle: lc, Topic: "schäfer", Channel: "schäfer", Config: cfg,
					Handler: handler, Marshaller: marshaller.NewMsgPack(),
					Version: version,
				},
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := &tnsq.Config{LookupHost: "invalid_host"}
		handler := test.NewHandler(nil)

		Convey("When I register a consumer", func() {
			err = tnsq.RegisterConsumer(
				tnsq.ConsumerParams{
					Lifecycle: lc, Topic: "topic", Channel: "channel", Config: cfg,
					Handler: handler, Marshaller: marshaller.NewMsgPack(),
					Version: version,
				},
				tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)

		err = tnsq.RegisterConsumer(
			tnsq.ConsumerParams{
				Lifecycle: lc, Topic: "topic", Channel: "channel", Config: cfg,
				Handler: handler, Marshaller: marshaller.NewMsgPack(),
				Version: version,
			},
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
		)
		So(err, ShouldBeNil)

		producer := tnsq.NewProducer(
			tnsq.ProducerParams{Lifecycle: lc, Config: cfg, Marshaller: marshaller.NewMsgPack(), Version: version},
			tnsq.WithProducerLogger(logger), tnsq.WithProducerTracer(tracer), tnsq.WithProducerRetry(), tnsq.WithProducerBreaker(),
		)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := message.New([]byte("test"))
			err = producer.Publish(context.Background(), "topic", message)
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

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)

		err = tnsq.RegisterConsumer(
			tnsq.ConsumerParams{
				Lifecycle: lc, Topic: "topic", Channel: "channel", Config: cfg,
				Handler: handler, Marshaller: marshaller.NewMsgPack(),
				Version: version,
			},
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
		)
		So(err, ShouldBeNil)

		producer, _ := nsq.NewProducer(cfg.Host, nsq.NewConfig())

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

// nolint:goerr113
func TestReceiveError(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		version := version.Version("1.0.0")

		tracer, err := opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Config: test.NewJaegerConfig(), Version: version})
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(errors.New("something went wrong"))

		err = tnsq.RegisterConsumer(
			tnsq.ConsumerParams{
				Lifecycle: lc, Topic: "topic", Channel: "channel", Config: cfg,
				Handler: handler, Marshaller: marshaller.NewMsgPack(),
				Version: version,
			},
			tnsq.WithConsumerLogger(logger), tnsq.WithConsumerTracer(tracer),
		)
		So(err, ShouldBeNil)

		producer := tnsq.NewProducer(
			tnsq.ProducerParams{Lifecycle: lc, Config: cfg, Marshaller: marshaller.NewMsgPack(), Version: version},
			tnsq.WithProducerLogger(logger), tnsq.WithProducerTracer(tracer), tnsq.WithProducerRetry(), tnsq.WithProducerBreaker(),
		)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := message.New([]byte("test"))
			err = producer.Publish(context.Background(), "topic", message)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
