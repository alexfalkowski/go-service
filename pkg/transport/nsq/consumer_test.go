package nsq_test

import (
	"sync"
	"testing"
	"time"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgNSQ "github.com/alexfalkowski/go-service/pkg/transport/nsq"
	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

type handler struct {
	m   *nsq.Message
	mux sync.Mutex
}

func (h *handler) Message() *nsq.Message {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.m
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	return nil
}

func TestConsumer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := pkgNSQ.NewConfig()

		Convey("When I register a consumer", func() {
			lc := fxtest.NewLifecycle(t)
			params := &pkgNSQ.ConsumerParams{
				SystemConfig: systemConfig,
				NSQConfig:    nsqConfig,
				Topic:        "topic",
				Channel:      "channel",
				Handler:      &handler{},
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
		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := pkgNSQ.NewConfig()
		handler := &handler{}
		params := &pkgNSQ.ConsumerParams{
			SystemConfig: systemConfig,
			NSQConfig:    nsqConfig,
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
