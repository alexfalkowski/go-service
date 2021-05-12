package nsq_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestProducer(t *testing.T) {
	Convey("Given I have all the configuration", t, func() {
		cfg, err := nsq.NewConfig()
		So(err, ShouldBeNil)

		Convey("When I register a producer", func() {
			lc := fxtest.NewLifecycle(t)

			logger, err := zap.NewLogger(lc, zap.NewConfig())
			So(err, ShouldBeNil)

			params := &nsq.ProducerParams{
				Config: cfg,
				Logger: logger,
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
