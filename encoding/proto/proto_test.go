package proto_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestEncoder(t *testing.T) {
	Convey("Given I have proto encoder", t, func() {
		encoder := proto.NewEncoder()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I encode the proto", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(bytes.Bytes())

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "CAE=")
			})
		})
	})

	Convey("Given I have proto encoder", t, func() {
		encoder := proto.NewEncoder()

		Convey("When I decode the proto", func() {
			b, err := base64.StdEncoding.DecodeString("CAE=")
			So(err, ShouldBeNil)

			var msg grpc_health_v1.HealthCheckResponse

			err = encoder.Decode(bytes.NewReader(b), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}
