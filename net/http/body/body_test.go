package body_test

import (
	"bufio"
	"net/http/httptest"
	"testing"
	"testing/synctest"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestReadAllHandlesNilBody(t *testing.T) {
	req := &http.Request{}

	data, bufferedBody, err := body.ReadAll(req)
	require.NoError(t, err)
	require.Empty(t, data)
	require.NotNil(t, bufferedBody)
	require.Equal(t, http.NoBody, req.Body)
}

func TestReadAllBuffersBody(t *testing.T) {
	req := &http.Request{Body: &test.TrackedBody{Reader: strings.NewReader("body")}}

	data, bufferedBody, err := body.ReadAll(req)
	require.NoError(t, err)
	require.Equal(t, []byte("body"), data)
	require.NotNil(t, bufferedBody)
}

func TestCloseSkipsEmptyBody(t *testing.T) {
	body.Close(nil)
	body.Close(http.NoBody)
}

func TestCloseClosesBody(t *testing.T) {
	trackedBody := &test.TrackedBody{Reader: strings.NewReader("body")}

	body.Close(trackedBody)
	require.True(t, trackedBody.Closed)
}

func TestNewLazyHandlerRejectsContentLengthOverLimit(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/test", strings.NewReader(strings.Repeat("a", 100)))
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler := body.NewLazyHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		require.Fail(t, "next handler should not be called")
	}), 64)

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, res.Code)
	test.RequireTrimmedResponseBody(t, res, "http: request entity too large")
}

func TestNewLazyHandlerSkipsEmptyBody(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()

	handler := body.NewLazyHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.NoBody, req.Body)
		res.WriteHeader(http.StatusOK)
	}), 64)

	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func TestNewLazyHandlerDoesNotEnforceLimitMidStream(t *testing.T) {
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/test", &test.UnknownLengthReader{Reader: strings.NewReader(strings.Repeat("a", 100))})
	require.NoError(t, err)
	res := httptest.NewRecorder()

	var data []byte
	var readErr error
	handler := body.NewLazyHandler(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		data, _, readErr = io.ReadAll(req.Body)
	}), 64)

	handler.ServeHTTP(res, req)

	require.NoError(t, readErr)
	require.Len(t, data, 100)
}

func TestNewLazyHandlerStreamsBodyIncrementally(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		pipeReader, pipeWriter := io.Pipe()

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/test", pipeReader)
		require.NoError(t, err)
		req.ContentLength = -1

		res := httptest.NewRecorder()
		lines := make(chan string, 2)
		done := make(chan struct{})

		handler := body.NewLazyHandler(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			reader := bufio.NewReader(req.Body)

			for range 2 {
				line, lineErr := reader.ReadString('\n')
				require.NoError(t, lineErr)
				lines <- line
			}

			close(done)
		}), 1024)

		go handler.ServeHTTP(res, req)

		_, err = pipeWriter.Write([]byte("one\n"))
		require.NoError(t, err)
		requireLine(t, lines, "one\n")

		_, err = pipeWriter.Write([]byte("two\n"))
		require.NoError(t, err)
		requireLine(t, lines, "two\n")

		require.NoError(t, pipeWriter.Close())
		requireDone(t, done)
	})
}

func requireLine(tb testing.TB, lines chan string, want string) {
	tb.Helper()

	select {
	case line := <-lines:
		require.Equal(tb, want, line)
	case <-time.After(2 * time.Second):
		tb.Fatal("timed out waiting for line")
	}
}

func requireDone(tb testing.TB, done chan struct{}) {
	tb.Helper()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		tb.Fatal("timed out waiting for handler to finish")
	}
}
