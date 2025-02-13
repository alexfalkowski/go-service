//nolint:varnamelen
package http_test

import (
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	h "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

func TestSendReceiveWithRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(http.DefaultTransport), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		world.RegisterEvents(t.Context())

		Convey("When I send an event", func() {
			ctx := world.EventsContext(t.Context())

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

			world.RequireStop()
		})
	})
}

func TestSendReceiveWithoutRoundTripper(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		world.RegisterEvents(t.Context())

		Convey("When I send an event", func() {
			ctx := world.EventsContext(t.Context())

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

			world.RequireStop()
		})
	})
}

func TestSendNotReceive(t *testing.T) {
	Convey("Given I have a http event receiver", t, func() {
		world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRoundTripper(&delRoundTripper{rt: http.DefaultTransport}), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		world.RegisterEvents(t.Context())

		Convey("When I send an event", func() {
			ctx := world.EventsContext(t.Context())

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

			world.RequireStop()
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
