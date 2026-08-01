package client_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestStreamRecvsValues(t *testing.T) {
	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	var greetings []string
	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			for {
				var res test.Response
				if err := stream.Recv(&res); err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}

				greetings = append(greetings, res.Greeting)
			}
		})

	require.NoError(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice"}, greetings)
}

func TestStreamGetIssuesGetRequest(t *testing.T) {
	var method string

	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		method = req.Method
		handler.ServeHTTP(res, req)
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	var greeting string
	err := c.StreamGet(t.Context(), server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			var res test.Response
			if err := stream.Recv(&res); err != nil {
				return err
			}

			greeting = res.Greeting
			return nil
		})

	require.NoError(t, err)
	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "Hello Bob", greeting)
}

func TestStreamPostPutPatchIssueBidiRequests(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*client.Client, context.Context, string, client.Options, client.RequestStreamHandler) error
	}{
		{name: "post", method: http.MethodPost, call: (*client.Client).StreamPost},
		{name: "put", method: http.MethodPut, call: (*client.Client).StreamPut},
		{name: "patch", method: http.MethodPatch, call: (*client.Client).StreamPatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var method string

			handler := contentstream.NewRequestHandler(
				test.Content, test.StreamEncoder,
				contentstream.Options{},
				func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
					req, err := stream.Recv()
					if err != nil {
						return err
					}

					return stream.Send(&test.Response{Greeting: "Hello " + req.Name})
				},
			)

			server := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				method = req.Method
				handler.ServeHTTP(res, req)
			}))
			server.EnableHTTP2 = true
			server.StartTLS()
			t.Cleanup(server.Close)

			c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(server.Client().Transport))

			var greeting string
			err := tt.call(c, t.Context(), server.URL, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON},
				func(_ context.Context, stream *client.RequestResponseStream) error {
					if err := stream.Send(&test.Request{Name: "Bob"}); err != nil {
						return err
					}

					var res test.Response
					if err := stream.Recv(&res); err != nil {
						return err
					}

					greeting = res.Greeting
					return nil
				})

			require.NoError(t, err)
			require.Equal(t, tt.method, method)
			require.Equal(t, "Hello Bob", greeting)
		})
	}
}

func TestStreamRequiresStreamRegistry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.NDJSON)
		_, _ = io.WriteString(res, "{}\n")
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, nil, test.Pool)

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON}, func(_ context.Context, stream *client.ResponseStream) error {
		var res test.Response
		return stream.Recv(&res)
	})

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestStreamRejectsUnsupportedResponseMedia(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.JSON)
		_, _ = io.WriteString(res, "{}")
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{}, func(_ context.Context, _ *client.ResponseStream) error {
		return nil
	})

	require.Error(t, err)
}

func TestStreamSurfacesErrorResponseBeforeCallingHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Error)
		res.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(res, "bad request")
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	called := false
	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{}, func(_ context.Context, _ *client.ResponseStream) error {
		called = true
		return nil
	})

	require.Error(t, err)
	require.EqualError(t, err, "bad request")
	require.Equal(t, http.StatusBadRequest, status.Code(err))
	require.False(t, called)
}

func TestRequestStreamInterleavesOverHTTP2(t *testing.T) {
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
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

	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(server.Client().Transport))

	var greetings []string
	err := c.RequestStream(ctx, http.MethodPost, server.URL, client.Options{ContentType: media.NDJSON, Accept: media.NDJSON},
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

func TestRequestStreamRejectsUnsupportedRequestMedia(t *testing.T) {
	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	err := c.RequestStream(t.Context(), http.MethodPost, "http://example.invalid", client.Options{ContentType: media.JSON},
		func(_ context.Context, _ *client.RequestResponseStream) error {
			return nil
		})

	require.Error(t, err)
}

func TestRequestStreamRequiresStreamRegistry(t *testing.T) {
	c := client.NewClient(test.Content, nil, test.Pool)

	err := c.RequestStream(t.Context(), http.MethodPost, "http://example.invalid", client.Options{},
		func(_ context.Context, _ *client.RequestResponseStream) error {
			return nil
		})

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestRequestStreamSurfacesErrorResponseViaRecv(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.Error)
		res.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(res, "bad request")
	}))
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(server.Client().Transport))

	sendErr := errors.New("unset")
	err := c.RequestStream(ctx, http.MethodPost, server.URL, client.Options{ContentType: media.NDJSON},
		func(_ context.Context, stream *client.RequestResponseStream) error {
			sendErr = stream.Send(&test.Request{Name: "Bob"})

			var res test.Response
			return stream.Recv(&res)
		})

	require.NoError(t, sendErr)
	require.Error(t, err)
	require.EqualError(t, err, "bad request")
	require.Equal(t, http.StatusBadRequest, status.Code(err))
}

func TestRequestStreamSendIsSticky(t *testing.T) {
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			for {
				if _, err := stream.Recv(); err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}
			}
		})
	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(server.Client().Transport))

	var firstErr, secondErr error
	err := c.RequestStream(ctx, http.MethodPost, server.URL, client.Options{ContentType: media.NDJSON},
		func(_ context.Context, stream *client.RequestResponseStream) error {
			firstErr = stream.Send(&test.Unencodable{})
			secondErr = stream.Send(&test.Unencodable{})
			return firstErr
		})

	require.Error(t, firstErr)
	require.Equal(t, firstErr, secondErr)
	require.Equal(t, err, firstErr)
}

func TestStreamRejectsValueOverCap(t *testing.T) {
	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "a greeting far longer than the configured cap"})
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithMaxResponseSize(8))

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			var res test.Response
			return stream.Recv(&res)
		})

	require.Error(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
}

func TestStreamRejectsValueOverCapWhenDecodeSucceedsInOneRead(t *testing.T) {
	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "a greeting far longer than the configured cap"})
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = func(r io.Reader) encodingstream.Decoder { return &test.SingleReadDecoder{R: r} }
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	c := client.NewClient(cont, sm, test.Pool, client.WithMaxResponseSize(8))

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			var res test.Response
			return stream.Recv(&res)
		})

	require.Error(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
}

func TestStreamRecvRecoversCapReaderErrorDespiteDecoder(t *testing.T) {
	tests := []struct {
		name     string
		greeting string
		maxSize  bytes.Size
		decoder  func(io.Reader) encodingstream.Decoder
	}{
		{
			name:     "decoder discards the error's type",
			greeting: "a greeting far longer than the configured cap",
			maxSize:  8,
			decoder:  func(r io.Reader) encodingstream.Decoder { return &test.OpaqueErrorDecoder{R: r} },
		},
		{
			name:     "decoder calls Read again after budget.Reader already latched an error",
			greeting: "ab",
			maxSize:  1,
			decoder:  func(r io.Reader) encodingstream.Decoder { return &test.TripleReadDecoder{R: r} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
				return stream.Send(&test.Response{Greeting: tt.greeting})
			})

			server := httptest.NewServer(handler)
			t.Cleanup(server.Close)

			sm := encodingstream.NewMap()
			codec := sm.Get("json")
			codec.Decoder = tt.decoder
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, test.Pool)

			c := client.NewClient(cont, sm, test.Pool, client.WithMaxResponseSize(tt.maxSize))

			err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
				func(_ context.Context, stream *client.ResponseStream) error {
					var res test.Response
					return stream.Recv(&res)
				})

			require.Error(t, err)
			require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
		})
	}
}

func TestStreamCapIsPerValueNotCumulative(t *testing.T) {
	const values = 10

	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		for range values {
			if err := stream.Send(&test.Response{Greeting: "hi"}); err != nil {
				return err
			}

			time.Sleep(3 * time.Millisecond)
		}

		return nil
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithMaxResponseSize(200))

	count := 0
	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			for {
				var res test.Response
				if err := stream.Recv(&res); err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}

				count++
			}
		})

	require.NoError(t, err)
	require.Equal(t, values, count)
}

func TestStreamRecvDoesNotFalsePositiveOnCoalescedReads(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.NDJSON)
		_, _ = res.Write([]byte("{\"Greeting\":\"Hello Bob\"}\n{\"Greeting\":\"Hello Alice\"}\n"))
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithMaxResponseSize(32))

	var greetings []string
	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			for {
				var res test.Response
				if err := stream.Recv(&res); err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}

				greetings = append(greetings, res.Greeting)
			}
		})

	require.NoError(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice"}, greetings)
}

func TestStreamRecvRejectsBufferedValueOverCap(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.NDJSON)
		_, _ = res.Write([]byte("{\"Greeting\":\"A\"}\n{\"Greeting\":\"a greeting far longer than the configured cap\"}\n"))
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithMaxResponseSize(20))

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
		func(_ context.Context, stream *client.ResponseStream) error {
			var first test.Response
			require.NoError(t, stream.Recv(&first))

			var second test.Response
			return stream.Recv(&second)
		})

	require.Error(t, err)
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(err))
}

func TestStreamClosesResponseBodyOnHandlerPanic(t *testing.T) {
	body := &test.TrackedBody{Reader: strings.NewReader("{\"Greeting\":\"hi\"}\n")}
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{content.TypeKey: []string{media.NDJSON}},
			Body:       body,
			Request:    req,
		}, nil
	})

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(rt))

	require.PanicsWithValue(t, "boom", func() {
		_ = c.Stream(t.Context(), http.MethodGet, "http://example.com", client.Options{Accept: media.NDJSON},
			func(_ context.Context, _ *client.ResponseStream) error {
				panic("boom")
			})
	})

	require.True(t, body.Closed)
}

func TestRequestStreamClosesPipeOnHandlerPanic(t *testing.T) {
	readDone := make(chan error, 1)

	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		go func() {
			_, _, err := io.ReadAll(req.Body)
			readDone <- err
		}()

		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
	})

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(rt))

	require.PanicsWithValue(t, "boom", func() {
		_ = c.RequestStream(t.Context(), http.MethodPost, "http://example.com", client.Options{ContentType: media.NDJSON},
			func(_ context.Context, _ *client.RequestResponseStream) error {
				panic("boom")
			})
	})

	select {
	case <-readDone:
	case <-time.After(2 * time.Second):
		t.Fatal("pipe reader never unblocked after handler panic")
	}
}

func TestStreamSurvivesSlowValuesWithConfiguredUnaryTimeout(t *testing.T) {
	handler := contentstream.NewHandler(test.Content, test.StreamEncoder, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithTimeout(10*time.Millisecond))

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	var greetings []string
	err := c.Stream(ctx, http.MethodGet, server.URL, client.Options{Accept: media.NDJSON}, func(_ context.Context, stream *client.ResponseStream) error {
		for {
			var res test.Response
			if err := stream.Recv(&res); err != nil {
				if stream.IsFinished(err) {
					return nil
				}

				return err
			}

			greetings = append(greetings, res.Greeting)
		}
	})

	require.NoError(t, err)
	require.Equal(t, []string{"Hello Bob", "Hello Alice"}, greetings)
}

func TestStreamReturnsEncodeErrorFromNewRequest(t *testing.T) {
	enc := encoding.NewMap()
	enc.Register("json", test.NewEncoder(test.ErrFailed))
	cont := content.NewContent(enc, test.Pool)

	c := client.NewClient(cont, test.StreamEncoder, test.Pool)

	err := c.Stream(t.Context(), http.MethodGet, "http://example.com", client.Options{ContentType: media.JSON, Request: &test.Request{Name: "Bob"}},
		func(_ context.Context, _ *client.ResponseStream) error {
			return nil
		})

	require.Error(t, err)
	require.ErrorIs(t, err, test.ErrFailed)
}

func TestStreamReturnsTransportError(t *testing.T) {
	rt := test.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, test.ErrFailed
	})

	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool, client.WithRoundTripper(rt))

	err := c.Stream(t.Context(), http.MethodGet, "http://example.com", client.Options{}, func(_ context.Context, _ *client.ResponseStream) error {
		return nil
	})

	require.Error(t, err)
	require.ErrorIs(t, err, test.ErrFailed)
}

func TestRequestStreamRejectsNilEncoderForRegisteredKind(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = nil
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	c := client.NewClient(cont, sm, test.Pool)

	err := c.RequestStream(t.Context(), http.MethodPost, "http://example.invalid", client.Options{ContentType: media.NDJSON},
		func(_ context.Context, _ *client.RequestResponseStream) error {
			return nil
		})

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestRequestStreamReturnsNewRequestError(t *testing.T) {
	c := client.NewClient(test.Content, test.StreamEncoder, test.Pool)

	err := c.RequestStream(t.Context(), "BAD METHOD", "http://example.com", client.Options{ContentType: media.NDJSON},
		func(_ context.Context, _ *client.RequestResponseStream) error {
			return nil
		})

	require.Error(t, err)
}

func TestRequestStreamFinishReturnsEncoderCloseError(t *testing.T) {
	rt := test.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		go func() { _, _, _ = io.ReadAll(req.Body) }()
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
	})

	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return test.CloseErrEncoder{} }
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	c := client.NewClient(cont, sm, test.Pool, client.WithRoundTripper(rt))

	err := c.RequestStream(t.Context(), http.MethodPost, "http://example.com", client.Options{ContentType: media.NDJSON},
		func(_ context.Context, _ *client.RequestResponseStream) error {
			return nil
		})

	require.ErrorIs(t, err, test.ErrFailed)
}

func TestNewResponseDecoderRejectsNilDecoderForRegisteredKind(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = nil
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(content.TypeKey, media.NDJSON)
		_, _ = io.WriteString(res, "{}\n")
	}))
	t.Cleanup(server.Close)

	c := client.NewClient(cont, sm, test.Pool)

	err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON}, func(_ context.Context, stream *client.ResponseStream) error {
		var res test.Response
		return stream.Recv(&res)
	})

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestStreamRecvBufferedLenFallback(t *testing.T) {
	tests := []struct {
		name    string
		decoder encodingstream.Decoder
	}{
		{name: "no buffered method", decoder: test.NoBufferedDecoder{}},
		{name: "wrong buffered type", decoder: test.WrongBufferedDecoder{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := encodingstream.NewMap()
			codec := sm.Get("json")
			codec.Decoder = func(io.Reader) encodingstream.Decoder { return tt.decoder }
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, test.Pool)

			server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.Header().Set(content.TypeKey, media.NDJSON)
				_, _ = io.WriteString(res, "{}\n")
			}))
			t.Cleanup(server.Close)

			c := client.NewClient(cont, sm, test.Pool)

			err := c.Stream(t.Context(), http.MethodGet, server.URL, client.Options{Accept: media.NDJSON},
				func(_ context.Context, stream *client.ResponseStream) error {
					var res test.Response
					return stream.Recv(&res)
				})

			require.NoError(t, err)
		})
	}
}
