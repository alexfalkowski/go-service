package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestSecure(t *testing.T) {
	Convey("Given I a secure client", t, func() {
		world := test.NewWorld(t, test.WithWorldSecure(), test.WithWorldTelemetry("prometheus"))
		world.Start()

		Convey("When I query github", func() {
			client := world.Client.NewHTTP()

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://github.com/alexfalkowski", http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have valid response", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

		world.Stop()
	})
}
