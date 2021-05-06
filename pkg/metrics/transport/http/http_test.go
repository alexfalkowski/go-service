package http_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/logger/zap"
	metricsHTTP "github.com/alexfalkowski/go-service/pkg/metrics/transport/http"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx/fxtest"
)

func TestHTTP(t *testing.T) {
	Convey("Given I register the metrics handler", t, func() {
		lc := fxtest.NewLifecycle(t)
		mux := pkgHTTP.NewMux()
		logger, err := zap.NewLogger(lc)
		So(err, ShouldBeNil)

		cfg := &config.Config{HTTPPort: "10002"}

		pkgHTTP.Register(lc, test.NewShutdowner(), mux, cfg, logger)

		err = metricsHTTP.Register(mux)
		So(err, ShouldBeNil)

		lc.RequireStart()

		Convey("When I query metrics", func() {
			client := &http.Client{Transport: pkgHTTP.NewRoundTripper(logger)}

			req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:10002/metrics", nil)
			So(err, ShouldBeNil)

			resp, err := client.Do(req)
			So(err, ShouldBeNil)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			So(err, ShouldBeNil)

			lc.RequireStop()

			Convey("Then I should have valid metrics", func() {
				So(string(body), ShouldContainSubstring, "go_info")
			})
		})
	})
}
