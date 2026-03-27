package mvc_test

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/stretchr/testify/require"
)

func TestStaticPathValueRejectsTraversal(t *testing.T) {
	mux := http.NewServeMux()
	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: slog.Default()}),
		FileSystem:  test.FileSystem,
		Layout:      test.Layout,
	})
	require.True(t, mvc.StaticPathValue("/{file...}", "file", "static"))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/robots.txt", http.NoBody)
	handler, _ := mux.Handler(req)
	req.SetPathValue("file", "../views/hello.tmpl")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
}
