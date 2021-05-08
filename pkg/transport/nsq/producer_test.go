package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		systemConfig := &config.Config{
			NSQLookupHost: "localhost:4161",
			NSQHost:       "localhost:4150",
		}
		nsqConfig := nsq.NewConfig()

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)

			logger, err := zap.NewLogger(lc, zap.NewConfig())
			So(err, ShouldBeNil)

			params := &nsq.ProducerParams{
				SystemConfig: systemConfig,
				NSQConfig:    nsqConfig,
				Logger:       logger,
			}
			_, err = nsq.NewProducer(lc, params)

			lc.RequireStart()

			Convey("Then I should not have an error", func() {
				So(err, ShouldBeNil)
			})

			lc.RequireStop()
		})
	})
}
