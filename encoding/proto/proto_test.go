package proto_test

import (
	"bytes"
	"testing"

	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/test"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestBinaryEncoder(t *testing.T) {
	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I encode and decode the proto", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var msg grpc_health_v1.HealthCheckResponse

			err = encoder.Decode(bytes, &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have a status", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
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

		Convey("When I decode the proto", func() {
			var msg grpc_health_v1.HealthCheckResponse

			err := encoder.Decode(bytes.NewBufferString("status:SERVING"), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
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

		Convey("When I decode the proto", func() {
			var msg grpc_health_v1.HealthCheckResponse

			err := encoder.Decode(bytes.NewBufferString(`{"status":"SERVING"}`), &msg)
			So(err, ShouldBeNil)

			Convey("Then I should have valid msg", func() {
				So(msg.GetStatus(), ShouldEqual, grpc_health_v1.HealthCheckResponse_SERVING)
			})
		})
	})
}

func TestErrEncoder(t *testing.T) {
	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		Convey("When I decode the proto with an erroneous reader", func() {
			var msg grpc_health_v1.HealthCheckResponse

			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have text encoder", t, func() {
		encoder := proto.NewText()

		Convey("When I decode the proto with an erroneous reader", func() {
			var msg grpc_health_v1.HealthCheckResponse

			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have json encoder", t, func() {
		encoder := proto.NewJSON()

		Convey("When I decode the proto with an erroneous reader", func() {
			var msg grpc_health_v1.HealthCheckResponse

			err := encoder.Decode(&test.ErrReaderCloser{}, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}

//nolint:funlen
func TestErrMessage(t *testing.T) {
	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		Convey("When I encode the proto with an erroneous message", func() {
			var msg test.ErrProto

			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have text encoder", t, func() {
		encoder := proto.NewText()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		Convey("When I encode the proto with an erroneous message", func() {
			var msg test.ErrProto

			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have json encoder", t, func() {
		encoder := proto.NewJSON()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		Convey("When I encode the proto with an erroneous message", func() {
			var msg test.ErrProto

			err := encoder.Encode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I decode the proto with a erroneous message", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var msg test.ErrProto

			err = encoder.Decode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})

	Convey("Given I have binary encoder", t, func() {
		encoder := proto.NewBinary()

		bytes := test.Pool.Get()
		defer test.Pool.Put(bytes)

		msg := &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}

		Convey("When I decode the proto with a wrong type", func() {
			err := encoder.Encode(bytes, msg)
			So(err, ShouldBeNil)

			var msg string

			err = encoder.Decode(bytes, &msg)

			Convey("Then I should have an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
