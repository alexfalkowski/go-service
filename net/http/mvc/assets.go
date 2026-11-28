package mvc

import (
	"crypto/sha256"
	"fmt"
	"path"
	"text/template"

	"github.com/alexfalkowski/go-service/v3/bytes"
	"github.com/alexfalkowski/go-service/v3/errors"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/url"
	"github.com/alexfalkowski/go-service/v3/strings"
)

// Asset describes one static asset. Set either File or Manifest.
type Asset struct {
	// Name is the relative URL path below the asset prefix and the key used by assetURL.
	// Its extension determines the response content type.
	Name string

	// File is the application filesystem path of an asset served as-is.
	File string

	// Manifest is the application filesystem path of a text/template manifest.
	// Its include function inserts a file's raw contents using a filesystem-relative path.
	Manifest string
}

// NewAssets loads the declared assets from fsys and fingerprints their final bytes with SHA-256.
// Provide the result through [ServerParams.Assets] to mount the routes at startup.
// Asset routes follow the transport's configured authentication and access policy by default.
// Pass [StaticOption] values to configure authentication and caching for all routes in the declaration.
// Prefix is a clean absolute URL path such as /assets. Prefix and asset names use literal URL paths
// that require no escaping. The constructor returns an error on duplicate names, invalid paths,
// unreadable sources, or manifest errors.
func NewAssets(fsys io.FS, prefix string, assets []Asset, opts ...StaticOption) (*Assets, error) {
	if !strings.HasPrefix(prefix, "/") || path.Clean(prefix) != prefix || !isLiteralAssetPath(prefix) {
		return nil, fmt.Errorf("mvc: invalid asset prefix %q", prefix)
	}

	files := make(map[string]staticAsset, len(assets))
	for _, asset := range assets {
		if asset.Name == "." || !io.ValidPath(asset.Name) || !isLiteralAssetPath(asset.Name) {
			return nil, fmt.Errorf("mvc: invalid asset name %q", asset.Name)
		}
		if _, exists := files[asset.Name]; exists {
			return nil, fmt.Errorf("mvc: duplicate asset %q", asset.Name)
		}

		content, err := readAsset(fsys, asset)
		if err != nil {
			return nil, fmt.Errorf("mvc: load asset %q: %w", asset.Name, err)
		}

		route := path.Join(prefix, asset.Name)
		fingerprint := fmt.Sprintf("%x", sha256.Sum256(content))
		files[asset.Name] = staticAsset{
			route:       route,
			url:         route + "?v=" + fingerprint,
			fingerprint: fingerprint,
			content:     content,
		}
	}

	return &Assets{files: files, options: options(opts...)}, nil
}

// Assets holds the static asset bodies, versioned URLs, and route options prepared at startup.
// Its contents are immutable and may be shared across servers and requests.
type Assets struct {
	files   map[string]staticAsset
	options *staticOptions
}

func (a *Assets) assetURL(name string) (string, error) {
	asset, ok := a.files[name]
	if !ok {
		return "", fmt.Errorf("mvc: unknown asset %q", name)
	}
	return asset.url, nil
}

type staticAsset struct {
	route       string
	url         string
	fingerprint string
	content     []byte
}

func (s *Server) registerAssets(assets *Assets) {
	for _, asset := range assets.files {
		handler := func(res http.ResponseWriter, req *http.Request) {
			query, err := url.ParseQuery(req.URL.RawQuery)
			versions, hasVersion := query["v"]
			matchesFingerprint := len(versions) == 1 && versions[0] == asset.fingerprint
			invalidVersion := hasVersion && !matchesFingerprint

			if err != nil || invalidVersion {
				res.Header().Set("Cache-Control", "no-store")
				res.WriteHeader(http.StatusNotFound)
				return
			}

			setStaticCacheControl(res, assets.options)
			writeStaticFile(res, asset.route, bytes.NewReader(asset.content), int64(len(asset.content)))
		}

		s.registerStatic(asset.route, handler, assets.options)
	}
}

func isLiteralAssetPath(name string) bool {
	return (&url.URL{Path: name}).EscapedPath() == name
}

func readAsset(fsys io.FS, asset Asset) ([]byte, error) {
	if strings.IsEmpty(asset.File) == strings.IsEmpty(asset.Manifest) {
		return nil, errors.New("expected exactly one file or manifest")
	}
	if !strings.IsEmpty(asset.File) {
		return io.ReadFile(fsys, asset.File)
	}

	manifest, err := io.ReadFile(fsys, asset.Manifest)
	if err != nil {
		return nil, err
	}

	include := func(name string) (string, error) {
		content, err := io.ReadFile(fsys, name)
		return string(content), err
	}
	tmpl, err := template.New(asset.Manifest).Funcs(template.FuncMap{"include": include}).Parse(string(manifest))
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, nil); err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}
