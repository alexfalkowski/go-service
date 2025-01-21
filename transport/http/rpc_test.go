//nolint:varnamelen
package http_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/test"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive
)

func TestRPCNoContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.NoContent)

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)

				_, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(status.Code(err), ShouldEqual, http.StatusOK)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRPCNoRequest(t *testing.T) {
	for _, mt := range []string{"gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.NoContent)

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)

				_, err := client.Invoke(context.Background(), nil)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRPCWithContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)

				resp, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.Greeting, ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestProtobufRPC(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.ProtobufSayHello)

			Convey("When I post data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[v1.SayHelloRequest, v1.SayHelloResponse](url, rpc.WithClientContentType("application/"+mt))

				resp, err := client.Invoke(context.Background(), &v1.SayHelloRequest{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.GetMessage(), ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestBadUnmarshalRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post data", func() {
				header := http.Header{}
				header.Set("Content-Type", "application/"+mt)

				res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString("a bad payload"))
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(body, ShouldNotBeBlank)
					So(res.StatusCode, ShouldEqual, 400)
				})

				world.RequireStop()
			})
		})
	}
}

func TestErrorRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.ErrorSayHello)

			Convey("When I post data", func() {
				header := http.Header{}
				header.Set("Content-Type", "application/"+mt)

				enc := test.Encoder.Get(mt)

				b := test.Pool.Get()
				defer test.Pool.Put(b)

				err := enc.Encode(b, test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				res, body, err := world.ResponseWithBody(context.Background(), "http", world.ServerHost(), http.MethodPost, "hello", header, b)
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(body, ShouldEqual, "rpc: ohh no")
					So(res.StatusCode, ShouldEqual, 503)
				})

				world.RequireStop()
			})
		})
	}
}

func TestAllowedRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)))
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post authenticated data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))

				resp, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(resp.Greeting, ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestDisallowedRPC(t *testing.T) {
	for _, mt := range []string{"application/json", "application/yaml", "application/yml", "application/toml", "application/gob", "test"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t,
				test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
				test.WithWorldToken(nil, test.NewVerifier("test")),
			)
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post authenticated data", func() {
				url := fmt.Sprintf("http://%s/hello", world.ServerHost())
				client := rpc.NewClient[test.Request, test.Response](url,
					rpc.WithClientContentType(mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))

				_, err := client.Invoke(context.Background(), &test.Request{Name: "Bob"})

				Convey("Then I should have an error", func() {
					So(status.IsError(err), ShouldBeTrue)
					So(status.Code(err), ShouldEqual, http.StatusUnauthorized)
					So(err.Error(), ShouldContainSubstring, "token: invalid match")
				})

				world.RequireStop()
			})
		})
	}
}
