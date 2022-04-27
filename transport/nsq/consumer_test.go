package nsq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	tnsq "github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		Convey("When I register a consumer", func() {
			params := tnsq.ConsumerParams{
				Lifecycle:  lc,
				Config:     cfg,
				Logger:     logger,
				Topic:      "topic",
				Channel:    "channel",
				Tracer:     tracer,
				Handler:    test.NewHandler(nil),
				Marshaller: marshaller.NewMsgPack(),
			}
			err := tnsq.RegisterConsumer(params)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestInvalidConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		Convey("When I register a consumer", func() {
			params := tnsq.ConsumerParams{
				Lifecycle:  lc,
				Config:     cfg,
				Logger:     logger,
				Topic:      "schäfer",
				Channel:    "schäfer",
				Tracer:     tracer,
				Handler:    test.NewHandler(nil),
				Marshaller: marshaller.NewMsgPack(),
			}
			err := tnsq.RegisterConsumer(params)

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

		tracer, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)
		cparams := tnsq.ConsumerParams{
			Lifecycle:  lc,
			Config:     cfg,
			Logger:     logger,
			Topic:      "topic",
			Channel:    "channel",
			Tracer:     tracer,
			Handler:    handler,
			Marshaller: marshaller.NewMsgPack(),
		}
		pparams := tnsq.ProducerParams{
			Lifecycle:  lc,
			Config:     cfg,
			Logger:     logger,
			Tracer:     tracer,
			Marshaller: marshaller.NewMsgPack(),
		}

		producer := tnsq.NewProducer(pparams)

		err = tnsq.RegisterConsumer(cparams)
		So(err, ShouldBeNil)

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

		tracer, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)
		cparams := tnsq.ConsumerParams{
			Lifecycle:  lc,
			Config:     cfg,
			Logger:     logger,
			Topic:      "topic",
			Channel:    "channel",
			Tracer:     tracer,
			Handler:    handler,
			Marshaller: marshaller.NewMsgPack(),
		}

		err = tnsq.RegisterConsumer(cparams)
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

		tracer, err := jaeger.NewTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(errors.New("something went wrong"))
		cparams := tnsq.ConsumerParams{
			Lifecycle:  lc,
			Config:     cfg,
			Logger:     logger,
			Topic:      "topic",
			Channel:    "channel",
			Tracer:     tracer,
			Handler:    handler,
			Marshaller: marshaller.NewMsgPack(),
		}
		pparams := tnsq.ProducerParams{
			Lifecycle:  lc,
			Config:     cfg,
			Logger:     logger,
			Tracer:     tracer,
			Marshaller: marshaller.NewMsgPack(),
		}

		producer := tnsq.NewProducer(pparams)

		err = tnsq.RegisterConsumer(cparams)
		So(err, ShouldBeNil)

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
