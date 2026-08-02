package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestRestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterHandlers("/hello", test.RestNoContent)

			url := world.NamedServerURL("http", "hello")
			err := world.Rest.Do(t.Context(), method, url, rest.Options{})
			require.NoError(t, err)
		})
	}
}

func TestRestRequestNoContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestNoContent)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			opts := rest.Options{
				ContentType: media.JSON,
				Request:     req,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
		})
	}
}

func TestRestError(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP(), test.WithWorldLoggerConfig("tint"))

			test.RegisterHandlers("/hello", test.RestError)

			url := world.NamedServerURL("http", "hello")
			err := world.Rest.Do(t.Context(), method, url, rest.Options{})
			require.Error(t, err)
		})
	}
}

func TestRestNotFound(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)

	res, body, err := world.GetBody(t.Context(), world.NamedServerURL("http", "missing"), header)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
	require.Equal(t, "text/error; charset=utf-8", res.Header.Get(http.ContentTypeKey))
	require.Equal(t, "http: not found", body)
}

func TestRestRequestError(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestError)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			opts := rest.Options{
				ContentType: media.JSON,
				Request:     req,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.Error(t, err)
		})
	}
}

func TestRestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodDelete, http.MethodGet} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

			test.RegisterHandlers("/hello", test.RestContent)

			url := world.NamedServerURL("http", "hello")
			resp := &test.Response{}
			opts := rest.Options{
				Response: resp,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello Bob", resp.Greeting)
		})
	}
}

func TestRestRequestWithContent(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

			test.RegisterRequestHandlers("/hello", test.RestRequestContent)

			url := world.NamedServerURL("http", "hello")
			req := &test.Request{Name: "test"}
			resp := &test.Response{}
			opts := rest.Options{
				ContentType: media.JSON,
				Request:     req,
				Response:    resp,
			}
			err := world.Rest.Do(t.Context(), method, url, opts)
			require.NoError(t, err)
			require.Equal(t, "Hello test", resp.Greeting)
		})
	}
}

func TestRestRequestUsesAcceptForResponse(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	test.RegisterRequestHandlers("/hello", test.RestRequestContent)

	body := test.Pool.Get()
	defer test.Pool.Put(body)
	require.NoError(t, test.Encoder.Get("json").Encode(body, &test.Request{Name: "Bob"}))
	header := http.Header{}
	header.Set(http.ContentTypeKey, media.JSON)
	header.Set(http.AcceptKey, media.YAML)

	res, responseBody, err := world.ResponseWithBody(
		t.Context(),
		world.NamedServerURL("http", "hello"),
		http.MethodPost,
		header,
		body,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, media.YAML, res.Header.Get(http.ContentTypeKey))
	var response test.Response
	require.NoError(t, test.Encoder.Get("yaml").Decode(strings.NewReader(responseBody), &response))
	require.Equal(t, "Hello Bob", response.Greeting)
}

func TestRestInvalidStatusCode(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldHTTP())

	test.RegisterHandlers("/hello", test.RestInvalidStatusCode)

	url := world.NamedServerURL("http", "hello")
	err := world.Rest.Get(t.Context(), url, rest.Options{})
	require.Error(t, err)

	err = world.Rest.Delete(t.Context(), url, rest.Options{})
	require.Error(t, err)

	test.RegisterRequestHandlers("/hello", test.RestRequestInvalidStatusCode)

	url = world.NamedServerURL("http", "hello")
	req := &test.Request{}
	opts := rest.Options{Request: req}

	err = world.Rest.Post(t.Context(), url, opts)
	require.Error(t, err)

	err = world.Rest.Put(t.Context(), url, opts)
	require.Error(t, err)

	err = world.Rest.Patch(t.Context(), url, opts)
	require.Error(t, err)
}

func TestRestStreamPostRecvAndSendOverRealServer(t *testing.T) {
	world := test.NewStartedWorld(t, test.WithWorldSecure(), test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())

	rest.StreamPost("/hello", func(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
		for {
			req, err := stream.Recv()
			if err != nil {
				if stream.IsFinished(err) {
					return nil
				}

				return err
			}

			if err := stream.Send(&test.Response{Greeting: "Hello " + req.Name}); err != nil {
				return err
			}
		}
	})

	url := world.PathServerURL("https", "hello")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var greetings []string
	err := world.Rest.StreamPost(ctx, url, rest.Options{ContentType: media.NDJSON, Accept: media.NDJSON},
		func(_ context.Context, stream *client.RequestResponseStream) error {
			for _, name := range []string{"Bob", "Alice"} {
				if err := stream.Send(&test.Request{Name: name}); err != nil {
					return err
				}

				var res test.Response
				if err := stream.Recv(&res); err != nil {
					return err
				}

				greetings = append(greetings, res.Greeting)
			}

			return nil
		})
	require.NoError(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice"}, greetings)
}

func TestRestStreamPostRefreshesReadDeadlineOverH2C(t *testing.T) {
	config := test.NewInsecureTransportConfig()
	config.HTTP.Timeout = 100 * time.Millisecond
	world := test.NewWorld(t, test.WithWorldTransportConfig(config), test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
	rest.Register(world.Router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{
		ReadTimeout: config.HTTP.GetReadTimeout(), WriteTimeout: config.HTTP.GetWriteTimeout(), MaxReceiveSize: config.HTTP.GetMaxReceiveSize(),
	})
	rest.StreamPost("/hello", echoStreamServer)
	world.Start()

	transport := http.Transport(nil)
	transport.Protocols.SetHTTP1(false)

	c := client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(transport))

	url := world.PathServerURL("http", "hello")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var greetings []string
	err := c.StreamPost(ctx, url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, staleClientAfterFourEchoes(&greetings))
	require.Error(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice", "Hello Carol", "Hello Dan"}, greetings)
}

func TestRestStreamPostSurvivesReceiveOnlyActivePhase(t *testing.T) {
	config := test.NewInsecureTransportConfig()
	config.HTTP.Timeout = 100 * time.Millisecond
	world := test.NewWorld(t, test.WithWorldTransportConfig(config), test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
	rest.Register(world.Router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{
		ReadTimeout: config.HTTP.GetReadTimeout(), WriteTimeout: config.HTTP.GetWriteTimeout(), MaxReceiveSize: config.HTTP.GetMaxReceiveSize(),
	})
	rest.StreamPost("/hello", receiveOnlyActivePhaseServer)
	world.Start()

	transport := http.Transport(nil)
	transport.Protocols.SetHTTP1(false)

	c := client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(transport))

	url := world.PathServerURL("http", "hello")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var greeting string
	err := c.StreamPost(ctx, url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, receiveOnlyActivePhaseClient(&greeting))
	require.NoError(t, err)
	require.Equal(t, "Hello Bob, Alice, Carol, Dan, Erin, Frank", greeting)
}

func TestRestStreamPostSurvivesSendOnlyActivePhase(t *testing.T) {
	config := test.NewInsecureTransportConfig()
	config.HTTP.Timeout = 100 * time.Millisecond
	world := test.NewWorld(t, test.WithWorldTransportConfig(config), test.WithWorldTelemetry("otlp"), test.WithWorldRest(), test.WithWorldHTTP())
	rest.Register(world.Router, test.UnaryContent, test.StreamContent, test.Pool, stream.Options{
		ReadTimeout: config.HTTP.GetReadTimeout(), WriteTimeout: config.HTTP.GetWriteTimeout(), MaxReceiveSize: config.HTTP.GetMaxReceiveSize(),
	})
	rest.StreamPost("/hello", sendOnlyActivePhaseServer)
	world.Start()

	transport := http.Transport(nil)
	transport.Protocols.SetHTTP1(false)

	c := client.NewClient(test.UnaryContent, test.StreamContent, test.Pool, client.WithRoundTripper(transport))

	url := world.PathServerURL("http", "hello")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var greetings []string
	err := c.StreamPost(ctx, url, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON}, sendOnlyActivePhaseClient(&greetings))
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two", "three", "four", "five", "six", "Hello Eve"}, greetings)
}

func echoStreamServer(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if stream.IsFinished(err) {
				return nil
			}

			return err
		}

		if err := stream.Send(&test.Response{Greeting: "Hello " + req.Name}); err != nil {
			return err
		}
	}
}

func staleClientAfterFourEchoes(greetings *[]string) client.RequestStreamHandler {
	return func(_ context.Context, stream *client.RequestResponseStream) error {
		for index, name := range []string{"Bob", "Alice", "Carol", "Dan"} {
			if index > 0 {
				time.Sleep(40 * time.Millisecond)
			}

			if err := stream.Send(&test.Request{Name: name}); err != nil {
				return err
			}

			var res test.Response
			if err := stream.Recv(&res); err != nil {
				return err
			}

			*greetings = append(*greetings, res.Greeting)
		}

		time.Sleep(120 * time.Millisecond)
		return stream.Send(&test.Request{Name: "Eve"})
	}
}

func receiveOnlyActivePhaseServer(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
	var names []string

	for range 6 {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		names = append(names, req.Name)
	}

	return stream.Send(&test.Response{Greeting: "Hello " + strings.Join(", ", names...)})
}

func receiveOnlyActivePhaseClient(greeting *string) client.RequestStreamHandler {
	return func(_ context.Context, stream *client.RequestResponseStream) error {
		for i, name := range []string{"Bob", "Alice", "Carol", "Dan", "Erin", "Frank"} {
			if i > 0 {
				time.Sleep(30 * time.Millisecond)
			}

			if err := stream.Send(&test.Request{Name: name}); err != nil {
				return err
			}
		}

		var res test.Response
		if err := stream.Recv(&res); err != nil {
			return err
		}

		*greeting = res.Greeting
		return nil
	}
}

func sendOnlyActivePhaseServer(_ context.Context, stream *stream.RequestStream[test.Request, test.Response]) error {
	for i, name := range []string{"one", "two", "three", "four", "five", "six"} {
		if i > 0 {
			time.Sleep(30 * time.Millisecond)
		}

		if err := stream.Send(&test.Response{Greeting: name}); err != nil {
			return err
		}
	}

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	return stream.Send(&test.Response{Greeting: "Hello " + req.Name})
}

func sendOnlyActivePhaseClient(greetings *[]string) client.RequestStreamHandler {
	return func(_ context.Context, stream *client.RequestResponseStream) error {
		for range 6 {
			var res test.Response
			if err := stream.Recv(&res); err != nil {
				return err
			}

			*greetings = append(*greetings, res.Greeting)
		}

		if err := stream.Send(&test.Request{Name: "Eve"}); err != nil {
			return err
		}

		var res test.Response
		if err := stream.Recv(&res); err != nil {
			return err
		}

		*greetings = append(*greetings, res.Greeting)
		return nil
	}
}
