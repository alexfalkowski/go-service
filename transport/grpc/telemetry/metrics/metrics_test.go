package metrics_test

import (
	"context"
	"errors"
	"io"
	"testing"

	me "github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/test"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/metadata"
)

func TestClientStream(t *testing.T) {
	Convey("Given I have a client stream", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewPrometheusMeter(lc)
		c := me.MustInt64Counter(m, "test_count", "test")
		h := me.MustFloat64Histogram(m, "test_hist", "testing")

		st := metrics.ClientStream{
			Options:  metric.WithAttributes(),
			Received: c, Sent: c, Handled: c,
			HandledHistogram: h,
			ClientStream:     &clientStream{err: io.EOF},
		}

		lc.RequireStart()

		Convey("When I try to receive a message with an EOF", func() {
			err := st.RecvMsg(nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError, io.EOF)
			})
		})

		lc.RequireStop()
	})

	Convey("Given I have a client stream", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewPrometheusMeter(lc)
		c := me.MustInt64Counter(m, "test_count", "test")
		h := me.MustFloat64Histogram(m, "test_hist", "testing")

		st := metrics.ClientStream{
			Options:  metric.WithAttributes(),
			Received: c, Sent: c, Handled: c,
			HandledHistogram: h,
			ClientStream:     &clientStream{err: errors.ErrUnsupported},
		}

		lc.RequireStart()

		Convey("When I try to receive a message with an error", func() {
			err := st.RecvMsg(nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError, errors.ErrUnsupported)
			})
		})

		lc.RequireStop()
	})
}

type clientStream struct {
	err error
}

func (c *clientStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (c *clientStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (c *clientStream) CloseSend() error {
	return nil
}

func (c *clientStream) Context() context.Context {
	return context.Background()
}

func (c *clientStream) SendMsg(m any) error {
	return c.err
}

func (c *clientStream) RecvMsg(m any) error {
	return c.err
}
