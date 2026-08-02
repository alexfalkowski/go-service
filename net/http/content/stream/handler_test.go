package stream_test

import (
	"bufio"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/compress"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/net/server"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestNewHandlerSendsValuesAndOptsOutOfGzip(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, media.NDJSON, res.Header().Get(http.ContentTypeKey))
	require.NotEmpty(t, res.Header().Get(compress.HeaderNoCompression))

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.Equal(t, "Hello Alice", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewHandlerGzipHandlerPassesThroughUncompressed(t *testing.T) {
	inner := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})
	handler := compress.GzipHandler(inner)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	req.Header.Set("Accept-Encoding", "gzip")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Empty(t, res.Header().Get("Content-Encoding"))
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, bufio.NewScanner(res.Body)))
}

func TestNewHandlerReturnsErrorBeforeFirstSendAsHTTPError(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
		return status.Error(http.StatusBadRequest, test.ErrFailed.Error())
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(http.ContentTypeKey))
}

func TestNewHandlerClosesEncoderBeforeFirstSend(t *testing.T) {
	encoder := &closeTracker{}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return &trackedEncoder{tracker: encoder} }
	sm.Register("json", codec)
	handler := contentstream.NewHandler(
		contentstream.NewContent(sm, test.Pool),
		contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
			return status.Error(http.StatusBadRequest, test.ErrFailed.Error())
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, 1, encoder.closes)
}

func TestNewHandlerReturnsCloseErrorBeforeFirstSendAsHTTPError(t *testing.T) {
	encoder := &closeTracker{err: test.ErrFailed}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return &trackedEncoder{tracker: encoder} }
	sm.Register("json", codec)
	handler := contentstream.NewHandler(
		contentstream.NewContent(sm, test.Pool),
		contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
			return nil
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, 1, encoder.closes)
}

func TestNewHandlerPreservesHandlerErrorWhenEncoderCloseFails(t *testing.T) {
	closeErr := errors.New("encoder close")
	encoder := &closeTracker{err: closeErr}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return &trackedEncoder{tracker: encoder} }
	sm.Register("json", codec)
	handler := contentstream.NewHandler(
		contentstream.NewContent(sm, test.Pool),
		contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
			return status.SafeError(http.StatusBadRequest, test.ErrFailed)
		},
	)

	req := httptest.NewRequestWithContext(status.WithRequestError(t.Context()), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, 1, encoder.closes)
	require.ErrorIs(t, status.RequestError(req.Context()), test.ErrFailed)
	require.ErrorIs(t, status.RequestError(req.Context()), closeErr)
}

func TestNewHandlerClosesEncoderBeforeFirstSendOnDrain(t *testing.T) {
	encoder := &closeTracker{}
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = func(io.Writer) encodingstream.Encoder { return &trackedEncoder{tracker: encoder} }
	sm.Register("json", codec)
	drain := make(chan struct{})
	handler := contentstream.NewHandler(
		contentstream.NewContent(sm, test.Pool),
		contentstream.Options{Drain: drain}, func(ctx context.Context, _ *contentstream.Stream[test.Response]) error {
			close(drain)
			<-ctx.Done()

			return ctx.Err()
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusServiceUnavailable, res.Code)
	require.Equal(t, 1, encoder.closes)
}

func TestNewHandlerFirstSendEncodeFailureDoesNotCommit(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Unencodable]) error {
		return stream.Send(&test.Unencodable{})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(http.ContentTypeKey))
}

func TestNewHandlerAbortsAfterCommit(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return test.ErrFailed
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
}

func TestNewHandlerEndsCleanlyAfterCommitOnDrain(t *testing.T) {
	drain := make(chan struct{})
	handler := contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain}, func(ctx context.Context, stream *contentstream.Stream[test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()

			return ctx.Err()
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewHandlerAbortsAfterCommitOnDrainForUnrelatedError(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)
	drain := make(chan struct{})
	handler := http.NewTelemetryHandler(contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain}, func(ctx context.Context, stream *contentstream.Stream[test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()

			return test.ErrFailed
		}), "http.server")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())

	spans := exporter.Spans()
	require.Len(t, spans, 1)
	require.Equal(t, tracer.StatusCodeError, spans[0].Status().Code)
	require.Equal(t, test.ErrFailed.Error(), spans[0].Status().Description)
	require.NotEmpty(t, spans[0].Events())
}

func TestNewHandlerEndsCleanlyAfterCommitOnDrainForCombinedError(t *testing.T) {
	drain := make(chan struct{})
	handler := contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain}, func(ctx context.Context, stream *contentstream.Stream[test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()

			return errors.Join(test.ErrFailed, ctx.Err())
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewHandlerDoesNotSendAfterDrain(t *testing.T) {
	drain := make(chan struct{})
	handler := contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain}, func(ctx context.Context, stream *contentstream.Stream[test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			close(drain)
			<-ctx.Done()

			return stream.Send(&test.Response{Greeting: "Hello Alice"})
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewHandlerRejectsSendWhenDrainStarts(t *testing.T) {
	previous := runtime.MaxProcs(1)
	t.Cleanup(func() { runtime.MaxProcs(previous) })
	drain := server.NewDrain()
	var sendErr error
	handler := contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain.Done()}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
			drain.Start()
			sendErr = stream.Send(&test.Response{Greeting: "sent after drain"})

			return sendErr
		},
	)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.ErrorIs(t, sendErr, contentstream.ErrDraining)
	require.Equal(t, http.StatusServiceUnavailable, res.Code)
}

func TestNewHandlerReturnsServiceUnavailableOnDrainBeforeCommit(t *testing.T) {
	drain := make(chan struct{})
	close(drain)
	called := false
	handler := contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{Drain: drain}, func(ctx context.Context, _ *contentstream.Stream[test.Response]) error {
			called = true
			<-ctx.Done()

			return ctx.Err()
		})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusServiceUnavailable, res.Code)
	require.False(t, called)
}

func TestNewHandlerAbortsAfterCommitRecordsTraceError(t *testing.T) {
	exporter := test.EnableIsolatedSpanExporter(t)
	handler := http.NewTelemetryHandler(contentstream.NewHandler(
		test.StreamContent,
		contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
			if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
				return err
			}

			return test.ErrFailed
		}), "http.server")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	spans := exporter.Spans()
	require.Len(t, spans, 1)
	require.Equal(t, tracer.StatusCodeError, spans[0].Status().Code)
	require.Equal(t, test.ErrFailed.Error(), spans[0].Status().Description)
	require.NotEmpty(t, spans[0].Events())
}

func TestNewHandlerRejectsUnsupportedMedia(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.JSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotAcceptable, res.Code)
}

func TestNewHandlerResolvesWildcardOrBrowserStyleAccept(t *testing.T) {
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
			handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
				return stream.Send(&test.Response{Greeting: "Hello Bob"})
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
			if tt.accept != "" {
				req.Header.Set(http.AcceptKey, tt.accept)
			}

			if tt.contentType != "" {
				req.Header.Set(http.ContentTypeKey, tt.contentType)
			}

			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, media.NDJSON, res.Header().Get(http.ContentTypeKey))
		})
	}
}

func TestNewHandlerSendChargesOneLimiterTokenPerMessage(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 2}

	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, 0, limiter.Remaining)
}

func TestNewHandlerSendAbortsAfterCommitWhenLimiterExhausted(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 1}

	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		if err := stream.Send(&test.Response{Greeting: "Hello Bob"}); err != nil {
			return err
		}

		return stream.Send(&test.Response{Greeting: "Hello Alice"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})

	scanner := bufio.NewScanner(res.Body)
	require.Equal(t, "Hello Bob", decodeNDJSONGreeting(t, scanner))
	require.False(t, scanner.Scan())
}

func TestNewHandlerSendPreCommitDenialMapsTo429(t *testing.T) {
	limiter := &test.CountingLimiter{Remaining: 0}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), limiter))
	res := httptest.NewRecorder()

	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "Hello Bob"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusTooManyRequests, res.Code)
	require.Equal(t, media.Error+"; charset=utf-8", res.Header().Get(http.ContentTypeKey))
}

func TestNewHandlerSendIsSticky(t *testing.T) {
	var firstErr, secondErr error

	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Unencodable]) error {
		firstErr = stream.Send(&test.Unencodable{})
		secondErr = stream.Send(&test.Unencodable{})
		return firstErr
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Error(t, firstErr)
	require.Equal(t, firstErr, secondErr)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewHandlerSendExtendDeadlineFailure(t *testing.T) {
	opts := contentstream.Options{WriteTimeout: time.MustParseDuration("1s")}
	handler := contentstream.NewHandler(test.StreamContent, opts, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewHandlerSendCommitFailure(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := &test.ErrResponseWriter{}

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})
}

func TestNewHandlerSendFlushFailure(t *testing.T) {
	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := &test.NoFlushResponseWriter{}

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(res, req)
	})
}

func TestNewHandlerSendLimiterErrorMapsTo500(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	req = req.WithContext(meta.WithLimiter(req.Context(), test.ErrorLimiter{}))
	res := httptest.NewRecorder()

	handler := contentstream.NewHandler(test.StreamContent, contentstream.Options{}, func(_ context.Context, stream *contentstream.Stream[test.Response]) error {
		return stream.Send(&test.Response{Greeting: "hi"})
	})

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestNewHandlerRejectsNilEncoderForRegisteredKind(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = nil
	sm.Register("json", codec)
	handler := contentstream.NewHandler(contentstream.NewContent(sm, test.Pool), contentstream.Options{}, func(_ context.Context, _ *contentstream.Stream[test.Response]) error {
		return nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", http.NoBody)
	req.Header.Set(http.AcceptKey, media.NDJSON)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotAcceptable, res.Code)
}

func decodeNDJSONGreeting(t *testing.T, scanner *bufio.Scanner) string {
	t.Helper()

	require.True(t, scanner.Scan())

	var res map[string]any
	require.NoError(t, json.Unmarshal(scanner.Bytes(), &res))

	greeting, _ := res["Greeting"].(string)
	return greeting
}
