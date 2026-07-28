package spike_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	snoop "github.com/felixge/httpsnoop"
	"github.com/klauspost/compress/gzhttp"
)

const writeTimeout = 300 * time.Millisecond

// chain approximates the real server stack: gzip outside the mux, httpsnoop
// wrapping (as otelhttp, the logger, and the recovery handler all do).
func chain(handler http.Handler) http.Handler {
	wrapped := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(snoop.Wrap(res, snoop.Hooks{}), req)
	})

	return gzhttp.GzipHandler(wrapped)
}

// slowStream writes 5 values 150ms apart, so the whole stream (~750ms) outlives
// the 300ms WriteTimeout. With extend=true it pushes the write deadline forward
// before each value, which is what PLAN.md B3 specifies.
func slowStream(extend bool) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set("Content-Type", "application/x-ndjson")
		controller := http.NewResponseController(res)

		for i := range 5 {
			if extend {
				if err := controller.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
					fmt.Fprintf(res, "{\"deadlineErr\":%q}\n", err.Error())
				}
			}

			fmt.Fprintf(res, "{\"n\":%d}\n", i)
			_ = controller.Flush()
			time.Sleep(150 * time.Millisecond)
		}
	})
}

func serve(tb testing.TB, handler http.Handler, http2 bool) *httptest.Server {
	tb.Helper()

	server := httptest.NewUnstartedServer(handler)
	server.Config.WriteTimeout = writeTimeout
	server.EnableHTTP2 = http2

	if http2 {
		server.StartTLS()
	} else {
		server.Start()
	}

	tb.Cleanup(server.Close)
	return server
}

func get(tb testing.TB, server *httptest.Server) (string, error) {
	tb.Helper()

	res, err := server.Client().Get(server.URL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	return string(data), err
}

func TestDeadlineExtension(t *testing.T) {
	for _, tc := range []struct {
		name   string
		extend bool
		http2  bool
	}{
		{"h1/no-extend", false, false},
		{"h1/extend", true, false},
		{"h2/no-extend", false, true},
		{"h2/extend", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := get(t, serve(t, chain(slowStream(tc.extend)), tc.http2))
			lines := 0
			for _, r := range body {
				if r == '\n' {
					lines++
				}
			}
			t.Logf("%-14s values=%d err=%v", tc.name, lines, err)
			if body != "" && lines > 0 && len(body) < 200 {
				t.Logf("%-14s body=%q", tc.name, body)
			}
		})
	}
}

// TestBidiOverHTTP2 is the premise the whole bidi design rests on: can the
// handler read the request body incrementally while writing responses, and can
// the client write a pipe-backed body while reading the response?
func TestBidiOverHTTP2(t *testing.T) {
	echo := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/x-ndjson")
		controller := http.NewResponseController(res)
		decoder := json.NewDecoder(req.Body)
		encoder := json.NewEncoder(res)

		for {
			var in map[string]int
			if err := decoder.Decode(&in); err != nil {
				return
			}
			_ = encoder.Encode(map[string]int{"echo": in["n"]})
			_ = controller.Flush()
		}
	})

	for _, http2 := range []bool{false, true} {
		label := "h1"
		if http2 {
			label = "h2"
		}

		t.Run(label, func(t *testing.T) {
			server := serve(t, echo, http2)

			ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
			defer cancel()

			reader, writer := io.Pipe()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL, reader)
			if err != nil {
				t.Fatal(err)
			}
			req.ContentLength = -1

			type result struct {
				res *http.Response
				err error
			}
			responses := make(chan result, 1)
			go func() {
				res, err := server.Client().Do(req)
				responses <- result{res, err}
			}()

			// Write the first value, then wait for its echo before writing the
			// second. If the transport buffers a direction, this deadlocks and the
			// context deadline ends the test.
			if _, err := writer.Write([]byte("{\"n\":1}\n")); err != nil {
				t.Fatalf("%s: first write: %v", label, err)
			}

			got := <-responses
			if got.err != nil {
				t.Logf("%-3s interleaved=false do-err=%v", label, got.err)
				_ = writer.Close()
				return
			}
			defer got.res.Body.Close()

			scanner := bufio.NewScanner(got.res.Body)
			if !scanner.Scan() {
				t.Logf("%-3s interleaved=false (no first echo) err=%v", label, scanner.Err())
				_ = writer.Close()
				return
			}
			first := scanner.Text()

			if _, err := writer.Write([]byte("{\"n\":2}\n")); err != nil {
				t.Logf("%-3s second write failed: %v", label, err)
				return
			}
			second := ""
			if scanner.Scan() {
				second = scanner.Text()
			}
			_ = writer.Close()

			t.Logf("%-3s interleaved=%t first=%q second=%q", label, second != "", first, second)
		})
	}
}

func TestNoCompressionOptOut(t *testing.T) {
	handler := gzhttp.GzipHandler(http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.Header().Set(gzhttp.HeaderNoCompression, "1")
		res.Header().Set("Content-Type", "application/x-ndjson")
		for i := range 200 {
			fmt.Fprintf(res, "{\"n\":%d}\n", i)
		}
	}))

	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	t.Logf("opt-out: Content-Encoding=%q NoCompression-header-leaked=%t",
		res.Header.Get("Content-Encoding"), res.Header.Get(gzhttp.HeaderNoCompression) != "")
}
