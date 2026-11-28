package mvc_test

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http/httptest"
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/meta"
	"github.com/alexfalkowski/go-service/v3/net/http/mvc"
	"github.com/stretchr/testify/require"
)

func TestNewServerAppliesAssetAuthenticationOptions(t *testing.T) {
	for _, scenario := range []struct {
		name            string
		options         []mvc.StaticOption
		unauthenticated bool
	}{
		{name: "default policy"},
		{name: "unauthenticated", options: []mvc.StaticOption{mvc.WithStaticUnauthenticated()}, unauthenticated: true},
		{name: "restore authentication", options: []mvc.StaticOption{mvc.WithStaticUnauthenticated(), mvc.WithStaticAuthenticated()}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			assets, err := mvc.NewAssets(test.FileSystem, "/assets", []mvc.Asset{
				{Name: "robots.txt", File: "static/robots.txt"},
			}, scenario.options...)
			require.NoError(t, err)

			mux := http.NewServeMux()
			policy := http.NewRoutePolicy()
			mvc.NewServer(mvc.ServerParams{
				Router:      http.NewRouter(mux, policy),
				FunctionMap: template.FuncMap{},
				FileSystem:  test.FileSystem,
				Pool:        test.Pool,
				Layout:      test.Layout,
				Assets:      assets,
			})

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/assets/robots.txt", http.NoBody)
			res := httptest.NewRecorder()
			mux.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			require.Equal(t, scenario.unauthenticated, policy.IsUnauthenticated(req))
			require.Empty(t, res.Header().Get("Cache-Control"))
		})
	}
}

func TestStaticAssetsServeFingerprintedCSSAndJavaScript(t *testing.T) {
	const cacheControl = "public, max-age=60"

	files := assetTestFiles()
	assets, err := mvc.NewAssets(files, "/assets", []mvc.Asset{
		{Name: "logs.js", Manifest: "scripts/logs.tmpl"},
		{Name: "logs.css", Manifest: "styles/logs.tmpl"},
		{Name: "print.css", File: "styles/base.css"},
	}, mvc.WithStaticUnauthenticated(), mvc.WithCacheControl(cacheControl))
	require.NoError(t, err)

	mux := http.NewServeMux()
	policy := http.NewRoutePolicy()
	functions := template.FuncMap{}
	server := mvc.NewServer(mvc.ServerParams{
		Router:      http.NewRouter(mux, policy),
		FunctionMap: functions,
		FileSystem:  files,
		Pool:        test.Pool,
		Layout:      mvc.NewLayout("layout.tmpl", "layout.tmpl"),
		Assets:      assets,
	})
	require.Contains(t, functions, "assetURL")
	view := server.NewFullView("page.tmpl")
	require.True(t, server.Get("/logs", func(_ context.Context) (*mvc.View, *test.Page, error) {
		return view, nil, nil
	}))

	// Requests should use the startup snapshot even if the source changes.
	files["styles/base.css"].Data[0] = 'B'

	get := func(ctx context.Context, url string) (*http.Request, *httptest.ResponseRecorder) {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)
		return req, res
	}

	_, page := get(t.Context(), "/logs")
	require.Equal(t, http.StatusOK, page.Code, "rendering the asset URLs")

	for _, scenario := range []struct {
		name        string
		content     string
		contentType string
	}{
		{"logs.js", "const label = '{{ready}}';\nconsole.log(label);\n", "javascript"},
		{"logs.css", "p { color: red; }\np { color: blue; }\n", "text/css"},
		{"print.css", "p { color: red; }\n", "text/css"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			url := fmt.Sprintf("/assets/%s?v=%x", scenario.name, sha256.Sum256([]byte(scenario.content)))
			test.RequireResponseBodyContains(t, page, `="`+url+`"`)

			req, res := get(t.Context(), url)
			require.Equal(t, http.StatusOK, res.Code)
			test.RequireResponseBody(t, res, scenario.content)
			require.Contains(t, res.Header().Get(http.ContentTypeKey), scenario.contentType)
			require.Equal(t, strconv.Itoa(len(scenario.content)), res.Header().Get(http.ContentLengthKey))
			require.Equal(t, cacheControl, res.Header().Get("Cache-Control"))
			require.True(t, policy.IsUnauthenticated(req))
		})
	}

	_, fragment := get(t.Context(), "/assets/setup.js")
	require.Equal(t, http.StatusNotFound, fragment.Code, "unlisted fragment")
}

func TestStaticAssetsRejectMismatchedVersions(t *testing.T) {
	const (
		assetPath    = "/assets/app.js"
		oldContent   = "console.log('old');"
		newContent   = "console.log('new');"
		cacheControl = "public, max-age=60"
	)

	newSnapshot := func(t *testing.T, content string) *http.ServeMux {
		t.Helper()
		files := fstest.MapFS{"app.js": &fstest.MapFile{Data: []byte(content)}}
		assets, err := mvc.NewAssets(files, "/assets", []mvc.Asset{
			{Name: "app.js", File: "app.js"},
		}, mvc.WithCacheControl(cacheControl))
		require.NoError(t, err)

		mux := http.NewServeMux()
		mvc.NewServer(mvc.ServerParams{
			Router:      newTestRouter(mux),
			FunctionMap: template.FuncMap{},
			FileSystem:  files,
			Pool:        test.Pool,
			Layout:      test.Layout,
			Assets:      assets,
		})
		return mux
	}

	old, current := newSnapshot(t, oldContent), newSnapshot(t, newContent)
	version := fmt.Sprintf("%x", sha256.Sum256([]byte(newContent)))
	versionedPath := assetPath + "?v=" + version
	for _, scenario := range []struct {
		name, path, cache, body string
		mux                     *http.ServeMux
		status                  int
	}{
		{"current fingerprint", versionedPath, cacheControl, newContent, current, http.StatusOK},
		{"new fingerprint on old server", versionedPath, "no-store", "", old, http.StatusNotFound},
		{"unversioned request", assetPath + "?other=1", cacheControl, oldContent, old, http.StatusOK},
		{"empty fingerprint", assetPath + "?v=", "no-store", "", current, http.StatusNotFound},
		{"malformed fingerprint", assetPath + "?v=%ZZ", "no-store", "", current, http.StatusNotFound},
		{"multiple fingerprints", versionedPath + "&v=" + version, "no-store", "", current, http.StatusNotFound},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, scenario.path, http.NoBody)
			res := httptest.NewRecorder()
			scenario.mux.ServeHTTP(res, req)

			require.Equal(t, scenario.status, res.Code)
			require.Equal(t, scenario.cache, res.Header().Get("Cache-Control"))
			test.RequireResponseBody(t, res, scenario.body)
		})
	}
}

func TestViewRejectsUnknownAsset(t *testing.T) {
	files := fstest.MapFS{
		"layout.tmpl": &fstest.MapFile{Data: []byte(`before {{assetURL "missing.js"}}`)},
		"page.tmpl":   &fstest.MapFile{},
		"app.js":      &fstest.MapFile{Data: []byte("console.log('loaded');")},
	}
	assets, err := mvc.NewAssets(files, "/assets", []mvc.Asset{{Name: "app.js", File: "app.js"}})
	require.NoError(t, err)

	server := mvc.NewServer(mvc.ServerParams{
		Router:      newTestRouter(http.NewServeMux()),
		FunctionMap: template.FuncMap{},
		FileSystem:  files,
		Pool:        test.Pool,
		Layout:      mvc.NewLayout("layout.tmpl", "layout.tmpl"),
		Assets:      assets,
	})
	view := server.NewFullView("page.tmpl")
	res := httptest.NewRecorder()
	ctx := meta.WithRequestResponse(t.Context(), nil, res)

	require.ErrorContains(t, view.Render(ctx, nil), `unknown asset "missing.js"`)
	test.RequireEmptyResponseBody(t, res)
}

func TestNewAssetsRejectsInvalidSources(t *testing.T) {
	files := fstest.MapFS{
		"missing.tmpl": &fstest.MapFile{Data: []byte(`{{include "missing.js"}}`)},
		"invalid.tmpl": &fstest.MapFile{Data: []byte(`{{include`)},
		"main.js":      &fstest.MapFile{Data: []byte("console.log('loaded');")},
	}
	t.Run("invalid prefix", func(t *testing.T) {
		assets, err := mvc.NewAssets(files, "/assets/", nil)
		require.ErrorContains(t, err, "invalid asset prefix")
		require.Nil(t, assets)
	})

	for _, scenario := range []struct {
		name    string
		message string
		sources []mvc.Asset
		err     error
	}{
		{"missing include", "missing.js", []mvc.Asset{{Name: "main.js", Manifest: "missing.tmpl"}}, io.ErrNotExist},
		{"missing manifest", "absent.tmpl", []mvc.Asset{{Name: "main.js", Manifest: "absent.tmpl"}}, io.ErrNotExist},
		{"malformed template", "invalid.tmpl", []mvc.Asset{{Name: "main.js", Manifest: "invalid.tmpl"}}, nil},
		{"missing file", "missing.css", []mvc.Asset{{Name: "main.css", File: "missing.css"}}, io.ErrNotExist},
		{"no source", "expected exactly one", []mvc.Asset{{Name: "main.js"}}, nil},
		{"two sources", "expected exactly one", []mvc.Asset{{Name: "main.js", File: "main.js", Manifest: "missing.tmpl"}}, nil},
		{"duplicate name", "duplicate asset", []mvc.Asset{{Name: "main.js", File: "main.js"}, {Name: "main.js", File: "main.js"}}, nil},
		{"route pattern", "invalid asset name", []mvc.Asset{{Name: "{file}", File: "main.js"}}, nil},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			assets, err := mvc.NewAssets(files, "/assets", scenario.sources)
			require.ErrorContains(t, err, scenario.message)
			require.Nil(t, assets)
			if scenario.err != nil {
				require.ErrorIs(t, err, scenario.err)
			}
		})
	}
}

func assetTestFiles() fstest.MapFS {
	return fstest.MapFS{
		"layout.tmpl": &fstest.MapFile{Data: []byte(`<link rel="stylesheet" href="{{assetURL "logs.css"}}">
<link rel="stylesheet" media="print" href="{{assetURL "print.css"}}">
<script src="{{assetURL "logs.js"}}"></script>`)},
		"page.tmpl":         &fstest.MapFile{},
		"scripts/logs.tmpl": &fstest.MapFile{Data: []byte(`{{include "scripts/setup.js"}}{{include "scripts/init.js"}}`)},
		"scripts/setup.js":  &fstest.MapFile{Data: []byte("const label = '{{ready}}';\n")},
		"scripts/init.js":   &fstest.MapFile{Data: []byte("console.log(label);\n")},
		"styles/logs.tmpl":  &fstest.MapFile{Data: []byte(`{{include "styles/base.css"}}{{include "styles/theme.css"}}`)},
		"styles/base.css":   &fstest.MapFile{Data: []byte("p { color: red; }\n")},
		"styles/theme.css":  &fstest.MapFile{Data: []byte("p { color: blue; }\n")},
	}
}
