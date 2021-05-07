package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgNSQ "github.com/alexfalkowski/go-service/pkg/transport/nsq"
	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

type handler struct{}

func (h *handler) HandleMessage(m *nsq.Message) error {
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
