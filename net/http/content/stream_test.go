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
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/klauspost/compress/gzhttp"
	"github.com/stretchr/testify/require"
)

func TestNewStreamHandlerSendsValuesAndOptsOutOfGzip(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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
	require.NotEmpty(t, res.Header().Get(gzhttp.HeaderNoCompression))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewStreamHandlerGzipHandlerPassesThroughUncompressed(t *testing.T) {
	inner := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})
	handler := gzhttp.GzipHandler(inner)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Empty(t, res.Header().Get("Content-Encoding"))
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, bufio.NewScanner(res.Body)))
}

func TestNewStreamHandlerReturnsErrorBeforeFirstSendAsHTTPError(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, _ *content.Stream[test.Response]) error {
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
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Unencodable]) error {
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
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, _ *content.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
}

func TestNewStreamHandlerSendChargesOneLimiterTokenPerMessage(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 2}

	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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

	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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

	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusTooManyRequests, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(content.TypeKey))
}

func TestNewRequestStreamHandlerRejectsHTTP1(t *testing.T) {
	handler := content.NewRequestStreamHandler(test.Content, 0, 0, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
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
	handler := content.NewRequestStreamHandler(test.Content, 0, 0, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
	handler := content.NewRequestStreamHandler(test.Content, 0, 0, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
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

			handler := content.NewRequestStreamHandler(test.Content, 0, tt.cap, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

func TestNewRequestStreamHandlerRecvRejectsValueOverCap(t *testing.T) {
	var recvErr error

	handler := content.NewRequestStreamHandler(test.Content, 0, 16, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

	handler := content.NewRequestStreamHandler(test.Content, 0, 16, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
			sm.RegisterDecoder("json", tt.decoder)
			cont := content.NewContent(test.Encoder, sm, test.Pool)

			var recvErr error

			handler := content.NewRequestStreamHandler(cont, 0, tt.maxSize, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

	handler := content.NewRequestStreamHandler(test.Content, 0, 24, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

	handler := content.NewRequestStreamHandler(test.Content, 0, 0, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
	handler := content.NewRequestStreamHandler(test.Content, 0, 16, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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

	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Unencodable]) error {
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
	handler := content.NewStreamHandler(test.Content, time.MustParseDuration("1s"), func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewStreamHandlerSendCommitFailure(t *testing.T) {
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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
	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
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

	handler := content.NewStreamHandler(test.Content, 0, func(_ context.Context, stream *content.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewRequestStreamHandlerRecvLimiterErrorMapsTo500(t *testing.T) {
	var recvErr error

	handler := content.NewRequestStreamHandler(test.Content, 0, 0, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
			sm.RegisterDecoder("json", func(io.Reader) stream.Decoder { return tt.decoder })
			cont := content.NewContent(test.Encoder, sm, test.Pool)

			handler := content.NewRequestStreamHandler(cont, 0, 0, func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
	sm.RegisterDecoder("json", func(r io.Reader) stream.Decoder { return &test.TripleReadDecoder{R: r} })
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	handler := content.NewRequestStreamHandler(cont, 0, bytes.Size(1), func(_ context.Context, stream *content.RequestStream[test.Request, test.Response]) error {
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
	sm.RegisterEncoder("json", nil)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	handler := content.NewStreamHandler(cont, 0, func(_ context.Context, _ *content.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
}

func TestNewRequestStreamHandlerRejectsNilDecoderForRegisteredKind(t *testing.T) {
	sm := stream.NewMap()
	sm.RegisterDecoder("json", nil)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	handler := content.NewRequestStreamHandler(cont, 0, 0, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
}

func TestNewRequestStreamHandlerRejectsNilEncoderForRegisteredKind(t *testing.T) {
	sm := stream.NewMap()
	sm.RegisterEncoder("json", nil)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	handler := content.NewRequestStreamHandler(cont, 0, 0, func(_ context.Context, _ *content.RequestStream[test.Request, test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/hello", http.NoBody)
	req.ProtoMajor = 2
	req.Header.Set(content.TypeKey, media.NDJSON)
	req.Header.Set(content.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Code)
}

func decodeNDJSONGreeting(t *testing.T, scanner *bufio.Scanner) string {
	t.Helper()

	require.True(t, scanner.Scan())

	var res map[string]any
	require.NoError(t, json.Unmarshal(scanner.Bytes(), &res))

	greeting, _ := res["Greeting"].(string)
	return greeting
}
