//nolint:varnamelen
package http_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/alexfalkowski/go-service/internal/test"
	v1 "github.com/alexfalkowski/go-service/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/net/http/status"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRPCNoContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.NoContent)

			Convey("When I post data", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)
				req := &test.Request{Name: "Bob"}
				res := &test.Response{}
				err := client.Invoke(t.Context(), "/hello", req, res)

				Convey("Then I should have no error", func() {
					So(err, ShouldBeNil)
					So(status.Code(err), ShouldEqual, http.StatusOK)
				})

				world.RequireStop()
			})
		})
	}
}

func TestRPCWithContent(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post data", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport),
					rpc.WithClientTimeout("10s"),
				)
				req := &test.Request{Name: "Bob"}
				res := &test.Response{}
				err := client.Invoke(t.Context(), "/hello", req, res)

				Convey("Then I should have response", func() {
					So(err, ShouldBeNil)
					So(res.Greeting, ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestSuccessProtobufRPC(t *testing.T) {
	for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessProtobufSayHello)

			Convey("When I post data", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url, rpc.WithClientContentType("application/"+mt))
				req := &v1.SayHelloRequest{Name: "Bob"}
				res := &v1.SayHelloResponse{}
				err := client.Invoke(t.Context(), "/hello", req, res)

				Convey("Then I should have response", func() {
					So(err, ShouldBeNil)
					So(res.GetMessage(), ShouldEqual, "Hello Bob")
				})

				world.RequireStop()
			})
		})
	}
}

func TestErroneousProtobufRPC(t *testing.T) {
	handlers := []content.RequestHandler[v1.SayHelloRequest, v1.SayHelloResponse]{
		test.ErrorsProtobufSayHello,
		test.ErrorsNotMappedProtobufSayHello,
	}

	for _, handler := range handlers {
		for _, mt := range []string{"proto", "protobuf", "prototext", "protojson"} {
			Convey("Given I have all the servers", t, func() {
				world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
				world.Register()
				world.RequireStart()

				rpc.Route("/hello", handler)

				Convey("When I post data", func() {
					url := "http://" + world.InsecureServerHost()
					client := rpc.NewClient(url, rpc.WithClientContentType("application/"+mt))
					req := &v1.SayHelloRequest{Name: "Bob"}
					res := &v1.SayHelloResponse{}
					err := client.Invoke(t.Context(), "/hello", req, res)

					Convey("Then I should have an error", func() {
						So(err, ShouldBeError)
						So(status.IsError(err), ShouldBeTrue)
						So(status.Code(err), ShouldEqual, http.StatusInternalServerError)
					})

					world.RequireStop()
				})
			})
		}
	}
}

func TestErroneousUnmarshalRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post data", func() {
				header := http.Header{}
				header.Set("Content-Type", "application/"+mt)

				res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, bytes.NewBufferString("an erroneous payload"))
				So(err, ShouldBeNil)

				Convey("Then I should have response", func() {
					So(body, ShouldNotBeBlank)
					So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
				})

				world.RequireStop()
			})
		})
	}
}

func TestErrorRPC(t *testing.T) {
	handlers := []content.RequestHandler[test.Request, test.Response]{
		test.ErrorSayHello,
		test.ErrorNotMappedSayHello,
	}

	for _, handler := range handlers {
		for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
			Convey("Given I have all the servers", t, func() {
				world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
				world.Register()
				world.RequireStart()

				rpc.Route("/hello", handler)

				Convey("When I post data", func() {
					header := http.Header{}
					header.Set("Content-Type", "application/"+mt)

					enc := test.Encoder.Get(mt)

					b := test.Pool.Get()
					defer test.Pool.Put(b)

					err := enc.Encode(b, test.Request{Name: "Bob"})
					So(err, ShouldBeNil)

					res, body, err := world.ResponseWithBody(t.Context(), "http", world.InsecureServerHost(), http.MethodPost, "hello", header, b)
					So(err, ShouldBeNil)

					Convey("Then I should have response", func() {
						So(body, ShouldEqual, "rpc: failed")
						So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
					})

					world.RequireStop()
				})
			})
		}
	}
}

func TestAllowedRPC(t *testing.T) {
	for _, mt := range []string{"json", "yaml", "yml", "toml", "gob"} {
		Convey("Given I have all the servers", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post authenticated data", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))
				req := &test.Request{Name: "Bob"}
				res := &test.Response{}
				err := client.Invoke(t.Context(), "/hello", req, res)

				Convey("Then I should have response", func() {
					So(err, ShouldBeNil)
					So(res.Greeting, ShouldEqual, "Hello Bob")
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
				test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldHTTP(),
			)
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I post authenticated data", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType(mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))
				req := &test.Request{Name: "Bob"}
				res := &test.Response{}
				err := client.Invoke(t.Context(), "/hello", req, res)

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

func TestInvalidRPCRequest(t *testing.T) {
	for _, mt := range []string{"gob"} {
		Convey("Given I have routes ready", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I try to send a hello request", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))
				err := client.Invoke(t.Context(), "/hello", nil, &test.Response{})

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}

func TestInvalidRPCResponse(t *testing.T) {
	for _, mt := range []string{"json"} {
		Convey("Given I have routes ready", t, func() {
			world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldHTTP())
			world.Register()
			world.RequireStart()

			rpc.Route("/hello", test.SuccessSayHello)

			Convey("When I try to send a hello request", func() {
				url := "http://" + world.InsecureServerHost()
				client := rpc.NewClient(url,
					rpc.WithClientContentType("application/"+mt),
					rpc.WithClientRoundTripper(world.NewHTTP().Transport))
				err := client.Invoke(t.Context(), "/hello", &test.Request{Name: "Bob"}, nil)

				Convey("Then I should have an error", func() {
					So(err, ShouldBeError)
				})

				world.RequireStop()
			})
		})
	}
}
