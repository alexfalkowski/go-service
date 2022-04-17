package nsq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/logger/zap"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		Convey("When I register a consumer", func() {
			params := &nsq.ConsumerParams{
				Lifecycle: lc,
				Config:    cfg,
				Logger:    logger,
				Topic:     "topic",
				Channel:   "channel",
				Tracer:    tracer,
				Handler:   test.NewHandler(nil),
			}
			err := nsq.RegisterConsumer(params)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
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

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(nil)
		cparams := &nsq.ConsumerParams{
			Lifecycle: lc,
			Config:    cfg,
			Logger:    logger,
			Topic:     "topic",
			Channel:   "channel",
			Tracer:    tracer,
			Handler:   handler,
		}
		pparams := &nsq.ProducerParams{
			Lifecycle: lc,
			Config:    cfg,
			Logger:    logger,
			Tracer:    tracer,
		}

		producer, err := nsq.NewProducer(pparams)
		So(err, ShouldBeNil)

		err = nsq.RegisterConsumer(cparams)
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

// nolint:goerr113
func TestReceiveError(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		tracer, err := opentracing.NewJaegerTransportTracer(lc, logger, test.NewJaegerConfig())
		So(err, ShouldBeNil)

		cfg := test.NewNSQConfig()
		handler := test.NewHandler(errors.New("something went wrong"))
		cparams := &nsq.ConsumerParams{
			Lifecycle: lc,
			Config:    cfg,
			Logger:    logger,
			Topic:     "topic",
			Channel:   "channel",
			Tracer:    tracer,
			Handler:   handler,
		}
		pparams := &nsq.ProducerParams{
			Lifecycle: lc,
			Config:    cfg,
			Logger:    logger,
			Tracer:    tracer,
		}

		producer, err := nsq.NewProducer(pparams)
		So(err, ShouldBeNil)

		err = nsq.RegisterConsumer(cparams)
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
