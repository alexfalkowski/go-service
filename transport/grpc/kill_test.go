package grpc_test

import (
	"context"
	"testing"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"go.uber.org/fx/fxtest"
)

func TestKill(t *testing.T) {
	Convey("Given I have a gRPC server", t, func() {
		lc := fxtest.NewLifecycle(t)
		logger := test.NewLogger(lc)
		cfg := test.NewInsecureTransportConfig()
		m := test.NewMeter(lc)
		hs := test.NewHTTPServer(lc, logger, test.NewOTLPTracerConfig(), cfg, m, nil)
		gs := test.NewGRPCServer(lc, logger, test.NewOTLPTracerConfig(), cfg, false, m, nil, nil)

		test.RegisterTransport(lc, gs, hs)
		lc.RequireStart()

		Convey("When I query for a greet with a dead server", func() {
			test.KillPort(cfg.GRPC.Port)

			ctx := meta.WithAttribute(context.Background(), "test", meta.Redacted("test"))

			conn := test.NewGRPCClient(lc, logger, cfg, test.NewOTLPTracerConfig(), nil, m)
			defer conn.Close()

			client := v1.NewGreeterServiceClient(conn)
			req := &v1.SayHelloRequest{Name: "test"}

			_, err := client.SayHello(ctx, req)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})

			lc.RequireStop()
		})
	})
}
