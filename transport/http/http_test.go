package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSecure(t *testing.T) {
	Convey("Given I a secure client", t, func() {
		world := test.NewWorld(t, test.WithWorldSecure(), test.WithWorldTelemetry("prometheus"), test.WithWorldHTTP())
		world.Register()
		world.RequireStart()

		Convey("When I query github", func() {
			client := world.NewHTTP()

			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://github.com/alexfalkowski", http.NoBody)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			Convey("Then I should have valid response", func() {
				So(resp.StatusCode, ShouldEqual, 200)
			})
		})

		world.RequireStop()
	})
}
