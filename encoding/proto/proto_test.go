package proto_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestBinaryEncoder(t *testing.T) {
	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		Convey("When I encode and decode the proto", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var decode grpc_health_v1.HealthCheckResponse

			err = encoder.Decode(bytes, &decode)
			So(err, ShouldBeNil)

			Convey("Then I should have a status", func() {
				So(decode.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})

		Convey("When I encode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			var msg string
			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I decode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var decode string
			err = encoder.Decode(bytes, &decode)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

//nolint:dupl
func TestTextEncoder(t *testing.T) {
	Convey("Given I have text encoder", t, func() {
		encoder := proto.NewText()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I encode the proto", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid proto", func() {
				So(bytes.String(), ShouldEqual, "status:SERVING")
			})
		})
	})

	Convey("Given I have text encoder", t, func() {
		encoder := proto.NewText()
		var msg grpc_health_v1.HealthCheckResponse

		Convey("When I decode the proto", func() {
			err := encoder.Decode(bytes.NewBufferString("status:SERVING"), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})

		Convey("When I encode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			var msg string
			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I decode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var decode string
			err = encoder.Decode(bytes, &decode)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

//nolint:dupl
func TestJSONEncoder(t *testing.T) {
	Convey("Given I have json encoder", t, func() {
		encoder := proto.NewJSON()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I encode the proto", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid proto", func() {
				So(bytes.String(), ShouldEqual, `{"status":"SERVING"}`)
			})
		})
	})

	Convey("Given I have json encoder", t, func() {
		encoder := proto.NewJSON()
		var msg grpc_health_v1.HealthCheckResponse

		Convey("When I decode the proto", func() {
			err := encoder.Decode(bytes.NewBufferString(`{"status":"SERVING"}`), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})

		Convey("When I encode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			var msg string
			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})

		Convey("When I decode with invalid type", func() {
			bytes := test.Pool.Get()
			defer test.Pool.Put(bytes)

			msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var decode string
			err = encoder.Decode(bytes, &decode)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

func TestErrEncoder(t *testing.T) {
	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()
		var msg grpc_health_v1.HealthCheckResponse

		Convey("When I decode the proto with an erroneous reader", func() {
			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have text encoder", t, func() {
		encoder := proto.NewText()
		var msg grpc_health_v1.HealthCheckResponse

		Convey("When I decode the proto with an erroneous reader", func() {
			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have json encoder", t, func() {
		encoder := proto.NewJSON()
		var msg grpc_health_v1.HealthCheckResponse

		Convey("When I decode the proto with an erroneous reader", func() {
			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
