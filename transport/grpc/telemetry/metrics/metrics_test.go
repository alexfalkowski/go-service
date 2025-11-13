package metrics_test

import (
	"io"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/metadata"
)

//nolint:dupl
func TestClientStream(t *testing.T) {
	Convey("Given I have a client stream", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewPrometheusMeter(lc)
		c := metrics.NewClient(test.Name, m)
		st := c.Stream(&clientStream{err: io.EOF}, metric.WithAttributes())

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
		c := metrics.NewClient(test.Name, m)
		st := c.Stream(&clientStream{err: test.ErrFailed}, metric.WithAttributes())

		lc.RequireStart()

		Convey("When I try to receive a message with an error", func() {
			err := st.RecvMsg(nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError, test.ErrFailed)
			})
		})

		lc.RequireStop()
	})
}

//nolint:dupl
func TestServerStream(t *testing.T) {
	Convey("Given I have a server stream", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewPrometheusMeter(lc)
		s := metrics.NewServer(test.Name, m)
		st := s.Stream(&serverStream{err: io.EOF}, metric.WithAttributes())

		lc.RequireStart()

		Convey("When I try to receive a message with an EOF", func() {
			err := st.RecvMsg(nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError, io.EOF)
			})
		})

		lc.RequireStop()
	})

	Convey("Given I have a server stream", t, func() {
		lc := fxtest.NewLifecycle(t)
		m := test.NewPrometheusMeter(lc)
		s := metrics.NewServer(test.Name, m)
		st := s.Stream(&serverStream{err: test.ErrFailed}, metric.WithAttributes())

		lc.RequireStart()

		Convey("When I try to receive a message with an error", func() {
			err := st.RecvMsg(nil)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError, test.ErrFailed)
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

func (c *clientStream) SendMsg(_ any) error {
	return c.err
}

func (c *clientStream) RecvMsg(_ any) error {
	return c.err
}

type serverStream struct {
	err error
}

func (s *serverStream) SetHeader(metadata.MD) error {
	return nil
}

func (s *serverStream) SendHeader(metadata.MD) error {
	return nil
}

func (s *serverStream) SetTrailer(metadata.MD) {
}

func (s *serverStream) Context() context.Context {
	return context.Background()
}

func (s *serverStream) SendMsg(_ any) error {
	return s.err
}

func (s *serverStream) RecvMsg(_ any) error {
	return s.err
}
