package stream_test

import (
	"bufio"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewRequestHandlerClosesCodecsBeforeAbortAfterCommitOnDrain(t *testing.T) {
	encoder := &closeTracker{}
	decoder := &closeTracker{}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return &trackedEncoder{tracker: encoder} }
	codec.Decoder = func(io.Reader) encodingstream.Decoder { return &trackedDecoder{tracker: decoder} }
	sm.Register("json", codec)
	drain := make(chan struct{})
	handler := contentstream.NewRequestHandler(
		test.Content, sm,
		contentstream.Options{Drain: drain}, func(ctx context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()

			return test.ErrFailed
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	require.Equal(t, 1, encoder.closes)
	require.Equal(t, 1, decoder.closes)
}

func TestNewRequestHandlerClosesDecoderWhenEncoderCloseFails(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)
	decoder := &closeTracker{err: errors.New("decoder close")}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return test.CloseErrEncoder{} }
	codec.Decoder = func(io.Reader) encodingstream.Decoder { return &trackedDecoder{tracker: decoder} }
	sm.Register("json", codec)
	handler := http.NewTelemetryHandler(contentstream.NewRequestHandler(
		test.Content, sm,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			return stream.Send(&test.Response{Greeting: "Hello Bob"})
		},
	), "http.server")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	require.Equal(t, 1, decoder.closes)
	spans := exporter.Spans()
	require.Len(t, spans, 1)
	require.Equal(t, test.ErrFailed.Error(), spans[0].Status().Description)
}

func TestNewRequestHandlerReturnsServiceUnavailableOnDrainBeforeHandler(t *testing.T) {
	drain := make(chan struct{})
	close(drain)
	called := false
	handler := contentstream.NewRequestHandler(
		test.Content, test.StreamEncoder,
		contentstream.Options{Drain: drain},
		func(_ context.Context, _ *contentstream.RequestStream[test.Request, test.Response]) error {
			called = true

			return nil
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusServiceUnavailable, res.Code)
	require.False(t, called)
}

func TestNewRequestHandlerUnblocksRecvOnDrain(t *testing.T) {
	drain := make(chan struct{})
	pipeReader, pipeWriter := io.Pipe()
	defer pipeWriter.Close()
	started := make(chan struct{})
	recvErr := make(chan error, 1)
	handler := contentstream.NewRequestHandler(
		test.Content, test.StreamEncoder,
		contentstream.Options{Drain: drain},
		func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			close(started)

			_, err := stream.Recv()
			recvErr <- err

			return err
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", pipeReader)
	req.ProtoMajor = 2
	req.ContentLength = -1
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()
	done := make(chan struct{})

	go func() {
		handler.ServeHTTP(res, req)
		close(done)
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for handler to start receiving")
	}

	close(drain)

	select {
	case <-done:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for drain to unblock receive")
	}

	require.ErrorIs(t, <-recvErr, contentstream.ErrDraining)
	require.Equal(t, http.StatusServiceUnavailable, res.Code)
}

func TestNewRequestHandlerRecvRejectsValueWhenDrainStartsAfterDecode(t *testing.T) {
	drain := make(chan struct{})
	body := &drainBody{closed: make(chan struct{})}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = func(io.Reader) encodingstream.Decoder {
		return &drainAfterDecode{drain: drain, closed: body.closed}
	}
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	var recvErr error
	handler := contentstream.NewRequestHandler(
		cont, sm,
		contentstream.Options{Drain: drain},
		func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()

			return recvErr
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", body)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.ErrorIs(t, recvErr, contentstream.ErrDraining)
	require.Equal(t, http.StatusServiceUnavailable, res.Code)
}

func TestNewRequestHandlerEndsCleanlyAfterCommitOnDrain(t *testing.T) {
	var recvErr error
	drain := make(chan struct{})
	handler := contentstream.NewRequestHandler(
		test.Content, test.StreamEncoder,
		contentstream.Options{Drain: drain},
		func(ctx context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()
			_, recvErr = stream.Recv()

			return recvErr
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.ErrorIs(t, recvErr, contentstream.ErrDraining)
	require.Equal(t, http.StatusOK, res.Code)
	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewRequestHandlerRejectsHTTP1(t *testing.T) {
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, _ *contentstream.RequestStream[test.Request, test.Response]) error {
			return nil
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 1
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusHTTPVersionNotSupported, res.Code)
}

func TestNewRequestHandlerRecvAndSend(t *testing.T) {
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

	body := "{\"Name\":\"Bob\"}\n{\"Name\":\"Alice\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.NDJSON, res.Header().Get(content.TypeKey))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewRequestHandlerRejectsUnsupportedRequestMedia(t *testing.T) {
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, _ *contentstream.RequestStream[test.Request, test.Response]) error {
			return nil
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.JSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
}

func TestNewRequestHandlerRecvUnderCapSucceeds(t *testing.T) {
	tests := []struct {
		name string
		cap  bytes.Size
	}{
		{name: "cap above the combined size of both values", cap: 64},
		{name: "cap between one value's size and the combined size of both", cap: 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var names []string

			opts := contentstream.Options{MaxReceiveSize: tt.cap}
			handler := contentstream.NewRequestHandler(
				test.Content,
				test.StreamEncoder,
				opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
					for {
						req, err := stream.Recv()
						if err != nil {
							if stream.IsFinished(err) {
								return nil
							}

							return err
						}

						names = append(names, req.Name)
					}
				})

			body := "{\"Name\":\"Bob\"}\n{\"Name\":\"Alice\"}\n"
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
			req.ProtoMajor = 2
			req.Header.Set(content.TypeKey, media.NDJSON)
			req.Header.Set(content.AcceptKey, media.NDJSON)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, []string{"Bob", "Alice"}, names)
		})
	}
}

func TestNewRequestHandlerRecvDeliversScalarValueAtExactCap(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()

	opts := contentstream.Options{MaxReceiveSize: bytes.Size(2)}
	handler := contentstream.NewRequestHandler(test.Content, test.StreamEncoder, opts, func(_ context.Context, stream *contentstream.RequestStream[int, int]) error {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		return stream.Send(req)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", pipeReader)
	req.ProtoMajor = 2
	req.ContentLength = -1
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		handler.ServeHTTP(res, req)
		close(done)
	}()

	_, err := pipeWriter.Write([]byte("42"))
	require.NoError(t, err)
	require.NoError(t, pipeWriter.Close())

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timed out waiting for handler to finish")
	}

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "42\n", res.Body.String())
}

func TestNewRequestHandlerRecvRejectsValueOverCap(t *testing.T) {
	var recvErr error

	opts := contentstream.Options{MaxReceiveSize: 16}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	body := "{\"Name\":\"a value with a name far longer than the configured cap\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)

	var maxBytesErr *http.MaxBytesError
	require.ErrorAs(t, recvErr, &maxBytesErr)
	require.Equal(t, int64(16), maxBytesErr.Limit)
	require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(recvErr))
}

func TestNewRequestHandlerRecvRejectsBufferedValueOverCap(t *testing.T) {
	var secondErr error

	opts := contentstream.Options{MaxReceiveSize: 16}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, err := stream.Recv()
			require.NoError(t, err)

			_, secondErr = stream.Recv()
			return secondErr
		})

	body := "{\"Name\":\"A\"}\n{\"Name\":\"a value with a name far longer than the configured cap\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)

	var maxBytesErr *http.MaxBytesError
	require.ErrorAs(t, secondErr, &maxBytesErr)
	require.Equal(t, int64(16), maxBytesErr.Limit)
}

func TestNewRequestHandlerRecvRejectsValueOverCapRegardlessOfDecoderBehavior(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		maxSize bytes.Size
		decoder func(io.Reader) encodingstream.Decoder
	}{
		{
			name:    "decode succeeds in one read",
			body:    "a value with far more than eight bytes",
			maxSize: 8,
			decoder: func(r io.Reader) encodingstream.Decoder { return &test.SingleReadDecoder{R: r} },
		},
		{
			name:    "decoder discards the error's type",
			body:    "{\"Name\":\"a value with a name far longer than the configured cap\"}\n",
			maxSize: 16,
			decoder: func(r io.Reader) encodingstream.Decoder { return &test.OpaqueErrorDecoder{R: r} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := encodingstream.NewMap()
			codec := sm.Get("json")
			codec.Decoder = tt.decoder
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, test.Pool)

			var recvErr error

			opts := contentstream.Options{MaxReceiveSize: tt.maxSize}
			handler := contentstream.NewRequestHandler(cont, sm, opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
				_, recvErr = stream.Recv()
				return recvErr
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(tt.body))
			req.ProtoMajor = 2
			req.Header.Set(content.TypeKey, media.NDJSON)
			req.Header.Set(content.AcceptKey, media.NDJSON)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)

			var maxBytesErr *http.MaxBytesError
			require.ErrorAs(t, recvErr, &maxBytesErr)
			require.Equal(t, int64(tt.maxSize), maxBytesErr.Limit)
			require.Equal(t, http.StatusRequestEntityTooLarge, status.Code(recvErr))
		})
	}
}

func TestNewRequestHandlerRecvCapIsPerValueNotCumulative(t *testing.T) {
	const values = 10

	var names []string

	pipeReader, pipeWriter := io.Pipe()

	opts := contentstream.Options{MaxReceiveSize: 24}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			for {
				req, err := stream.Recv()
				if err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}

				names = append(names, req.Name)
			}
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", pipeReader)
	req.ProtoMajor = 2
	req.ContentLength = -1
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		handler.ServeHTTP(res, req)
		close(done)
	}()

	for range values {
		_, err := pipeWriter.Write([]byte("{\"Name\":\"Bob\"}\n"))
		require.NoError(t, err)
	}
	require.NoError(t, pipeWriter.Close())

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		require.FailNow(t, "timed out waiting for handler to finish")
	}

	require.Equal(t, http.StatusOK, res.Code)
	require.Len(t, names, values)
}

func TestNewRequestHandlerRecvChargesOneLimiterTokenPerMessage(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 2}

	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			for {
				_, err := stream.Recv()
				if err != nil {
					if stream.IsFinished(err) {
						return nil
					}

					return err
				}
			}
		})

	body := "{\"Name\":\"Bob\"}\n{\"Name\":\"Alice\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, 0, limiter.Remaining)
}

func TestNewRequestHandlerRecvDoesNotChargeLimiterForOverCapValue(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 1}

	var recvErr error
	opts := contentstream.Options{MaxReceiveSize: 16}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	body := "{\"Name\":\"a value with a name far longer than the configured cap\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
	require.Equal(t, 1, limiter.Remaining)
}

func TestNewRequestHandlerSendExtendReadDeadlineFailure(t *testing.T) {
	opts := contentstream.Options{ReadTimeout: time.MustParseDuration("1s"), WriteTimeout: time.MustParseDuration("1s")}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			return stream.Send(&test.Response{Greeting: "hi"})
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.DeadlineResponseWriter{SetReadDeadlineFunc: func() error { return test.ErrFailed }}

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewRequestHandlerRecvFirstCallIsBoundedByReadTimeout(t *testing.T) {
	done := make(chan struct{})
	t.Cleanup(func() { <-done })

	pipeReader, pipeWriter := io.Pipe()
	t.Cleanup(func() { _ = pipeWriter.Close() })

	recvErr := make(chan error, 1)

	opts := contentstream.Options{ReadTimeout: time.MustParseDuration("200ms")}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, err := stream.Recv()
			recvErr <- err

			return err
		})

	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.StartTLS()
	t.Cleanup(server.Close)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL, pipeReader)
	require.NoError(t, err)
	req.ContentLength = -1
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)

	go func() {
		defer close(done)

		resp, doErr := server.Client().Do(req)
		if doErr == nil {
			_ = resp.Body.Close()
		}
	}()

	select {
	case err := <-recvErr:
		require.Error(t, err)
	case <-time.After(2 * time.Second):
		require.FailNow(t, "Recv did not return within the configured read timeout")
	}
}

func TestNewRequestHandlerRecvExtendReadDeadlineFailureBeforeDecode(t *testing.T) {
	var recvErr error

	opts := contentstream.Options{ReadTimeout: time.MustParseDuration("1s")}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("{\"Name\":\"Bob\"}\n"))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.DeadlineResponseWriter{SetReadDeadlineFunc: func() error { return test.ErrFailed }}

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.ErrorIs(t, recvErr, test.ErrFailed)
}

func TestNewRequestHandlerRecvExtendReadDeadlineFailureAfterDecode(t *testing.T) {
	var recvErr error
	calls := 0

	opts := contentstream.Options{ReadTimeout: time.MustParseDuration("1s")}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("{\"Name\":\"Bob\"}\n"))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.DeadlineResponseWriter{SetReadDeadlineFunc: func() error {
		calls++
		if calls == 1 {
			return nil
		}

		return test.ErrFailed
	}}

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.ErrorIs(t, recvErr, test.ErrFailed)
}

func TestNewRequestHandlerRecvExtendWriteDeadlineFailure(t *testing.T) {
	var recvErr error

	opts := contentstream.Options{ReadTimeout: time.MustParseDuration("1s"), WriteTimeout: time.MustParseDuration("1s")}
	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("{\"Name\":\"Bob\"}\n"))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.DeadlineResponseWriter{SetWriteDeadlineFunc: func() error { return test.ErrFailed }}

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.ErrorIs(t, recvErr, test.ErrFailed)
}

func TestNewRequestHandlerRecvLimiterErrorMapsTo500(t *testing.T) {
	var recvErr error

	handler := contentstream.NewRequestHandler(
		test.Content,
		test.StreamEncoder,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
			_, recvErr = stream.Recv()
			return recvErr
		})

	body := "{\"Name\":\"Bob\"}\n"
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader(body))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), test.ErrorLimiter{}))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, http.StatusInternalServerError, status.Code(recvErr))
}

func TestNewRequestHandlerRecvBufferedLenFallback(t *testing.T) {
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

			handler := contentstream.NewRequestHandler(cont, sm, contentstream.Options{}, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
				_, err := stream.Recv()
				return err
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("{}"))
			req.ProtoMajor = 2
			req.Header.Set(content.TypeKey, media.NDJSON)
			req.Header.Set(content.AcceptKey, media.NDJSON)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
		})
	}
}

func TestNewRequestHandlerRecvRecoversStickyCapReaderError(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = func(r io.Reader) encodingstream.Decoder { return &test.TripleReadDecoder{R: r} }
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, test.Pool)

	opts := contentstream.Options{MaxReceiveSize: bytes.Size(1)}
	handler := contentstream.NewRequestHandler(cont, sm, opts, func(_ context.Context, stream *contentstream.RequestStream[test.Request, test.Response]) error {
		_, err := stream.Recv()
		return err
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", strings.NewReader("ab"))
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
}

func TestNewRequestHandlerRejectsNilCodecFieldForRegisteredKind(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*encodingstream.Codec)
		status int
	}{
		{name: "nil decoder", mutate: func(c *encodingstream.Codec) { c.Decoder = nil }, status: http.StatusUnsupportedMediaType},
		{name: "nil encoder", mutate: func(c *encodingstream.Codec) { c.Encoder = nil }, status: http.StatusNotAcceptable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := encodingstream.NewMap()
			codec := sm.Get("json")
			tt.mutate(&codec)
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, test.Pool)

			handler := contentstream.NewRequestHandler(cont, sm, contentstream.Options{}, func(_ context.Context, _ *contentstream.RequestStream[test.Request, test.Response]) error {
				return nil
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
			req.ProtoMajor = 2
			req.Header.Set(content.TypeKey, media.NDJSON)
			req.Header.Set(content.AcceptKey, media.NDJSON)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, tt.status, res.Code)
		})
	}
}

type drainAfterDecode struct {
	drain  chan struct{}
	closed <-chan struct{}
}

func (d *drainAfterDecode) Decode(any) error {
	close(d.drain)
	<-d.closed

	return nil
}

func (*drainAfterDecode) Close() error {
	return nil
}

type drainBody struct {
	closed chan struct{}
}

func (*drainBody) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (b *drainBody) Close() error {
	close(b.closed)

	return nil
}

type closeTracker struct {
	closes int
	err    error
}

type trackedEncoder struct {
	tracker *closeTracker
}

func (*trackedEncoder) Encode(any) error {
	return nil
}

func (e *trackedEncoder) Close() error {
	e.tracker.closes++

	return nil
}

type trackedDecoder struct {
	tracker *closeTracker
}

func (*trackedDecoder) Decode(any) error {
	return nil
}

func (d *trackedDecoder) Close() error {
	d.tracker.closes++

	return d.tracker.err
}
