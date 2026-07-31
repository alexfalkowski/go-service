package content_test

import (
	"bufio"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/compress"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewStreamHandlerSendsValuesAndOptsOutOfGzip(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.NDJSON, res.Header().Get(content.TypeKey))
	require.NotEmpty(t, res.Header().Get(compress.HeaderNoCompression))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewStreamHandlerGzipHandlerPassesThroughUncompressed(t *testing.T) {
	inner := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})
	handler := compress.GzipHandler(inner)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Empty(t, res.Header().Get("Content-Encoding"))
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, bufio.NewScanner(res.Body)))
}

func TestNewStreamHandlerReturnsErrorBeforeFirstSendAsHTTPError(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, _ *content.Stream[test.Response]) error {
		return status.Error(http.StatusBadRequest, test.ErrFailed.Error())
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(content.TypeKey))
}

func TestNewStreamHandlerFirstSendEncodeFailureDoesNotCommit(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Unencodable]) error {
		return stream.Send(&test.Unencodable{})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(content.TypeKey))
}

func TestNewStreamHandlerAbortsAfterCommit(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return test.ErrFailed
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
}

func TestNewStreamHandlerRejectsUnsupportedMedia(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, _ *content.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotAcceptable, res.Code)
}

func TestNewStreamHandlerResolvesWildcardOrBrowserStyleAccept(t *testing.T) {
	tests := []struct {
		name        string
		accept      string
		contentType string
	}{
		{name: "curl default", accept: "*/*"},
		{name: "subtype wildcard", accept: "application/*"},
		{name: "browser-style list", accept: "text/html,application/xhtml+xml,*/*;q=0.8"},
		{name: "absent accept and content-type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
				return stream.Send(&test.Response{Greeting: "Hello Bob"})
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
			if tt.accept != "" {
				req.Header.Set(content.AcceptKey, tt.accept)
			}

			if tt.contentType != "" {
				req.Header.Set(content.TypeKey, tt.contentType)
			}

			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, media.NDJSON, res.Header().Get(content.TypeKey))
		})
	}
}

func TestNewStreamHandlerSendChargesOneLimiterTokenPerMessage(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 2}

	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, 0, limiter.Remaining)
}

func TestNewStreamHandlerSendAbortsAfterCommitWhenLimiterExhausted(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 1}

	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewStreamHandlerSendPreCommitDenialMapsTo429(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 0}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusTooManyRequests, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(content.TypeKey))
}

func TestNewRequestStreamHandlerRejectsHTTP1(t *testing.T) {
	handler := content.NewRequestStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvAndSend(t *testing.T) {
	handler := content.NewRequestStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRejectsUnsupportedRequestMedia(t *testing.T) {
	handler := content.NewRequestStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvUnderCapSucceeds(t *testing.T) {
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

			opts := content.StreamOptions{MaxReceiveSize: tt.cap}
			handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvDeliversScalarValueAtExactCap(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()

	opts := content.StreamOptions{MaxReceiveSize: bytes.Size(2)}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[int, int]) error {
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
		t.Fatal("timed out waiting for handler to finish")
	}

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "42\n", res.Body.String())
}

func TestNewRequestStreamHandlerRecvRejectsValueOverCap(t *testing.T) {
	var recvErr error

	opts := content.StreamOptions{MaxReceiveSize: 16}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvRejectsBufferedValueOverCap(t *testing.T) {
	var secondErr error

	opts := content.StreamOptions{MaxReceiveSize: 16}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvRejectsValueOverCapRegardlessOfDecoderBehavior(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		maxSize bytes.Size
		decoder func(io.Reader) stream.Decoder
	}{
		{
			name:    "decode succeeds in one read",
			body:    "a value with far more than eight bytes",
			maxSize: 8,
			decoder: func(r io.Reader) stream.Decoder { return &test.SingleReadDecoder{R: r} },
		},
		{
			name:    "decoder discards the error's type",
			body:    "{\"Name\":\"a value with a name far longer than the configured cap\"}\n",
			maxSize: 16,
			decoder: func(r io.Reader) stream.Decoder { return &test.OpaqueErrorDecoder{R: r} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := stream.NewMap()
			codec := sm.Get("json")
			codec.Decoder = tt.decoder
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, sm, test.Pool)

			var recvErr error

			opts := content.StreamOptions{MaxReceiveSize: tt.maxSize}
			handler := content.NewRequestStreamHandler(cont, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvCapIsPerValueNotCumulative(t *testing.T) {
	const values = 10

	var names []string

	pipeReader, pipeWriter := io.Pipe()

	opts := content.StreamOptions{MaxReceiveSize: 24}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
		t.Fatal("timed out waiting for handler to finish")
	}

	require.Equal(t, http.StatusOK, res.Code)
	require.Len(t, names, values)
}

func TestNewRequestStreamHandlerRecvChargesOneLimiterTokenPerMessage(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 2}

	handler := content.NewRequestStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvDoesNotChargeLimiterForOverCapValue(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 1}

	var recvErr error
	opts := content.StreamOptions{MaxReceiveSize: 16}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewStreamHandlerSendIsSticky(t *testing.T) {
	var firstErr, secondErr error

	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Unencodable]) error {
		firstErr = stream.Send(&test.Unencodable{})
		secondErr = stream.Send(&test.Unencodable{})
		return firstErr
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Error(t, firstErr)
	require.Equal(t, firstErr, secondErr)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewStreamHandlerSendExtendDeadlineFailure(t *testing.T) {
	opts := content.StreamOptions{WriteTimeout: time.MustParseDuration("1s")}
	handler := content.NewStreamHandler(test.Content, opts, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewRequestStreamHandlerSendExtendReadDeadlineFailure(t *testing.T) {
	opts := content.StreamOptions{ReadTimeout: time.MustParseDuration("1s"), WriteTimeout: time.MustParseDuration("1s")}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvFirstCallIsBoundedByReadTimeout(t *testing.T) {
	done := make(chan struct{})
	t.Cleanup(func() { <-done })

	pipeReader, pipeWriter := io.Pipe()
	t.Cleanup(func() { _ = pipeWriter.Close() })

	recvErr := make(chan error, 1)

	opts := content.StreamOptions{ReadTimeout: time.MustParseDuration("200ms")}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvExtendReadDeadlineFailureBeforeDecode(t *testing.T) {
	var recvErr error

	opts := content.StreamOptions{ReadTimeout: time.MustParseDuration("1s")}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvExtendReadDeadlineFailureAfterDecode(t *testing.T) {
	var recvErr error
	calls := 0

	opts := content.StreamOptions{ReadTimeout: time.MustParseDuration("1s")}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvExtendWriteDeadlineFailure(t *testing.T) {
	var recvErr error

	opts := content.StreamOptions{ReadTimeout: time.MustParseDuration("1s"), WriteTimeout: time.MustParseDuration("1s")}
	handler := content.NewRequestStreamHandler(test.Content, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewStreamHandlerSendCommitFailure(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.ErrResponseWriter{}

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})
}

func TestNewStreamHandlerSendFlushFailure(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := &test.NoFlushResponseWriter{}

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})
}

func TestNewStreamHandlerSendLimiterErrorMapsTo500(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), test.ErrorLimiter{}))
	res := httptest.NewRecorder()

	handler := content.NewStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewRequestStreamHandlerRecvLimiterErrorMapsTo500(t *testing.T) {
	var recvErr error

	handler := content.NewRequestStreamHandler(test.Content, content.StreamOptions{}, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvBufferedLenFallback(t *testing.T) {
	tests := []struct {
		name    string
		decoder stream.Decoder
	}{
		{name: "no buffered method", decoder: test.NoBufferedDecoder{}},
		{name: "wrong buffered type", decoder: test.WrongBufferedDecoder{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := stream.NewMap()
			codec := sm.Get("json")
			codec.Decoder = func(io.Reader) stream.Decoder { return tt.decoder }
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, sm, test.Pool)

			handler := content.NewRequestStreamHandler(cont, content.StreamOptions{}, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvRecoversStickyCapReaderError(t *testing.T) {
	sm := stream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = func(r io.Reader) stream.Decoder { return &test.TripleReadDecoder{R: r} }
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	opts := content.StreamOptions{MaxReceiveSize: bytes.Size(1)}
	handler := content.NewRequestStreamHandler(cont, opts, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewStreamHandlerRejectsNilEncoderForRegisteredKind(t *testing.T) {
	sm := stream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = nil
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	handler := content.NewStreamHandler(cont, content.StreamOptions{}, func(_ context.Context, _ *content.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotAcceptable, res.Code)
}

func TestNewRequestStreamHandlerRejectsNilCodecFieldForRegisteredKind(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*stream.Codec)
		status int
	}{
		{name: "nil decoder", mutate: func(c *stream.Codec) { c.Decoder = nil }, status: http.StatusUnsupportedMediaType},
		{name: "nil encoder", mutate: func(c *stream.Codec) { c.Encoder = nil }, status: http.StatusNotAcceptable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := stream.NewMap()
			codec := sm.Get("json")
			tt.mutate(&codec)
			sm.Register("json", codec)
			cont := content.NewContent(test.Encoder, sm, test.Pool)

			handler := content.NewRequestStreamHandler(cont, content.StreamOptions{}, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
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

func decodeNDJSONGreeting(t *testing.T, scanner *bufio.Scanner) string {
	t.Helper()

	require.True(t, scanner.Scan())

	var res map[string]any
	require.NoError(t, json.Unmarshal(scanner.Bytes(), &res))

	greeting, _ := res["Greeting"].(string)
	return greeting
}
