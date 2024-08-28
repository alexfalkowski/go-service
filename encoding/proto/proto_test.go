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

func TestMarshaller(t *testing.T) {
	Convey("Given I have proto marshaller", t, func() {
		m := proto.NewMarshaller()
		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I marshall the proto", func() {
			b, err := m.Marshal(msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(b)

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "CAE=")
			})
		})
	})

	Convey("Given I have proto marshaller", t, func() {
		m := proto.NewMarshaller()

		Convey("When I unmarshal the proto", func() {
			b, err := base64.StdEncoding.DecodeString("CAE=")
			So(err, ShouldBeNil)

			var msg grpc_health_v1.HealthCheckResponse

			err = m.Unmarshal(b, &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestEncoder(t *testing.T) {
	Convey("Given I have proto encoder", t, func() {
		e := proto.NewEncoder()

		b := test.Pool.Get()
		defer test.Pool.Put(b)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I encode the proto", func() {
			err := e.Encode(b, msg)
			So(err, ShouldBeNil)

			s := base64.StdEncoding.EncodeToString(b.Bytes())

			Convey("Then I should have valid proto", func() {
				So(s, ShouldEqual, "CAE=")
			})
		})
	})

	Convey("Given I have proto encoder", t, func() {
		e := proto.NewEncoder()

		Convey("When I decode the proto", func() {
			b, err := base64.StdEncoding.DecodeString("CAE=")
			So(err, ShouldBeNil)

			var msg grpc_health_v1.HealthCheckResponse

			err = e.Decode(bytes.NewReader(b), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}
