package nsq_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	pkgNSQ "github.com/alexfalkowski/go-service/pkg/transport/nsq"
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
		nsqConfig := pkgNSQ.NewConfig()

		Convey("When I register a consumer", func() {
			params := &pkgNSQ.ConsumerParams{
				SystemConfig: systemConfig,
				NSQConfig:    nsqConfig,
				Logger:       logger,
				Topic:        "topic",
				Channel:      "channel",
				Handler:      test.NewHandler(nil),
			}
			err := pkgNSQ.RegisterConsumer(lc, params)

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
		nsqConfig := pkgNSQ.NewConfig()
		handler := test.NewHandler(nil)
		params := &pkgNSQ.ConsumerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
			Topic:        "topic",
			Channel:      "channel",
			Handler:      handler,
		}

		producer, err := pkgNSQ.NewProducer(lc, systemConfig, nsqConfig)
		So(err, ShouldBeNil)

		err = pkgNSQ.RegisterConsumer(lc, params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			err = producer.Publish("topic", []byte("test"))
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
		nsqConfig := pkgNSQ.NewConfig()
		handler := test.NewHandler(errors.New("something went wrong"))
		params := &pkgNSQ.ConsumerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
			Logger:       logger,
			Topic:        "topic",
			Channel:      "channel",
			Handler:      handler,
		}

		producer, err := pkgNSQ.NewProducer(lc, systemConfig, nsqConfig)
		So(err, ShouldBeNil)

		err = pkgNSQ.RegisterConsumer(lc, params)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I send a message", func() {
			err = producer.Publish("topic", []byte("test"))
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			Convey("Then I should receive a message", func() {
				So(handler.Message(), ShouldNotBeNil)
			})

			lc.RequireStop()
		})
	})
}
