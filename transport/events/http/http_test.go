//nolint:varnamelen
package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

func TestSendReceiveWithRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport))
		world.Start()

		ctx := context.Background()

		world.RegisterEvents(ctx)

		Convey("When I send an event", func() {
			ctx := world.EventsContext(ctx)

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err := e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

			result := world.Sender.Send(ctx, e)

			Convey("Then I should receive an event", func() {
				So(protocol.IsACK(result), ShouldBeTrue)
				So(world.Event, ShouldNotBeNil)
				So(string(e.Data()), ShouldEqual, "test")
			})

			world.Stop()
		})
	})
}

func TestSendReceiveWithoutRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
		world.Start()

		ctx := context.Background()

		world.RegisterEvents(ctx)

		Convey("When I send an event", func() {
			ctx := world.EventsContext(ctx)

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err := e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

			result := world.Sender.Send(ctx, e)

			Convey("Then I should receive an event", func() {
				So(protocol.IsACK(result), ShouldBeTrue)
				So(world.Event, ShouldNotBeNil)
				So(string(e.Data()), ShouldEqual, "test")
			})

			world.Stop()
		})
	})
}

func TestSendNotReceive(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(&delRoundTripper{rt: http.DefaultTransport}))
		world.Start()

		ctx := context.Background()

		world.RegisterEvents(ctx)

		Convey("When I send an event", func() {
			ctx := world.EventsContext(ctx)

			e := events.NewEvent()
			e.SetSource("example/uri")
			e.SetType("example.type")

			err := e.SetData(events.TextPlain, "test")
			So(err, ShouldBeNil)

			result := world.Sender.Send(ctx, e)

			Convey("Then I should not receive an event", func() {
				So(protocol.IsNACK(result), ShouldBeTrue)
				So(world.Event, ShouldBeNil)
			})

			world.Stop()
		})
	})
}

type delRoundTripper struct {
	rt http.RoundTripper
}

func (r *delRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Del(h.HeaderWebhookID)

	return r.rt.RoundTrip(req)
}
