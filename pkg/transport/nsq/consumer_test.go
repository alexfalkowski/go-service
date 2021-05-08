package nsq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := nsq.NewConfig()

		Convey("When I register a consumer", func() {
			params := &nsq.ConsumerParams{
				SystemConfig: systemConfig,
				NSQConfig:    nsqConfig,
				Logger:       logger,
				Topic:        "topic",
				Channel:      "channel",
				Handler:      test.NewHandler(nil),
			}
			err := nsq.RegisterConsumer(lc, params)

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

		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := nsq.NewConfig()
		handler := test.NewHandler(nil)
		consumerParams := &nsq.ConsumerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
			Topic:        "topic",
			Channel:      "channel",
			Handler:      handler,
		}
		producerParams := &nsq.ProducerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
		}

		producer, err := nsq.NewProducer(lc, producerParams)
		So(err, ShouldBeNil)

		err = nsq.RegisterConsumer(lc, consumerParams)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := message.New([]byte("test"))
			_, err = producer.Publish(context.Background(), "topic", message)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}

func TestReceiveError(t *testing.T) {
	Convey("Given I have a consumer and a producer", t, func() {
		lc := fxtest.NewLifecycle(t)

		logger, err := zap.NewLogger(lc, zap.NewConfig())
		So(err, ShouldBeNil)

		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := nsq.NewConfig()
		handler := test.NewHandler(errors.New("something went wrong"))
		consumerParams := &nsq.ConsumerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
			Topic:        "topic",
			Channel:      "channel",
			Handler:      handler,
		}
		producerParams := &nsq.ProducerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
		}

		producer, err := nsq.NewProducer(lc, producerParams)
		So(err, ShouldBeNil)

		err = nsq.RegisterConsumer(lc, consumerParams)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			message := message.New([]byte("test"))
			_, err = producer.Publish(context.Background(), "topic", message)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
