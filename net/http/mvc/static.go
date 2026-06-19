package mvc

import (
	"io/fs"
	"path"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
)

// StaticFile registers an HTTP GET route that serves the named file from the registered filesystem.
//
// It returns false when MVC is not defined (see IsDefined).
func StaticFile(pattern, name string, opts ...StaticOption) bool {
	if !IsDefined() {
		return false
	}

	options := options(opts...)
	handler := func(res http.ResponseWriter, req *http.Request) {
		serveFile(res, req, name, options)
	}

	http.HandleFunc(mux, strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

// StaticPathValue registers an HTTP GET route that serves a file chosen by a path value.
//
// The file name is built under prefix from a validated request path value. Invalid paths and
// traversal attempts are rejected with HTTP 400.
//
// It returns false when MVC is not defined (see IsDefined).
func StaticPathValue(pattern, value, prefix string, opts ...StaticOption) bool {
	if !IsDefined() {
		return false
	}

	options := options(opts...)
	handler := func(res http.ResponseWriter, req *http.Request) {
		cleaned := path.Clean(req.PathValue(value))
		if cleaned == "." || cleaned != req.PathValue(value) || !fs.ValidPath(cleaned) || strings.Contains(cleaned, `\`) {
			res.WriteHeader(staticStatusCode(status.BadRequestError(fs.ErrInvalid)))
			return
		}

		name := path.Join(prefix, cleaned)
		serveFile(res, req, name, options)
	}

	http.HandleFunc(mux, strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

func serveFile(res http.ResponseWriter, req *http.Request, name string, options *staticOptions) {
	f, err := fileSystem.Open(name)
	if err != nil {
		res.WriteHeader(staticStatusCode(err))
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		res.WriteHeader(staticStatusCode(err))
		return
	}
	if info.IsDir() {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	setStaticCacheControl(res, options)
	if options.cacheValidators {
		serveFileWithCacheValidators(res, req, name, f, info)
		return
	}

	writeStaticFile(res, name, f, info)
}

func serveFileWithCacheValidators(res http.ResponseWriter, req *http.Request, name string, f fs.File, info fs.FileInfo) {
	etag := staticETag(name, info)
	setStaticValidators(res, etag, info)
	if staticETagMatches(req.Header.Get("If-None-Match"), etag) {
		res.WriteHeader(http.StatusNotModified)
		return
	}
	if strings.IsEmpty(req.Header.Get("If-None-Match")) && staticModifiedSinceMatches(req, info) {
		res.WriteHeader(http.StatusNotModified)
		return
	}

	writeStaticFile(res, name, f, info)
}

func writeStaticFile(res http.ResponseWriter, name string, f fs.File, info fs.FileInfo) {
	setStaticContentLength(res, info.Size())
	setStaticContentType(res, name)
	res.WriteHeader(http.StatusOK)
	_, _ = io.Copy(res, f)
}

func staticStatusCode(err error) int {
	if errors.Is(err, fs.ErrNotExist) {
		return http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return http.StatusForbidden
	}
	return status.Code(err)
}

func setStaticContentType(res http.ResponseWriter, name string) {
	mediaType := media.TypeByExtension(path.Ext(name))
	if !strings.IsEmpty(mediaType) {
		res.Header().Set(content.TypeKey, media.MustParse(mediaType).WithUTF8())
	}
}

func setStaticCacheControl(res http.ResponseWriter, options *staticOptions) {
	if !strings.IsEmpty(options.cacheControl) {
		res.Header().Set("Cache-Control", options.cacheControl)
	}
}

func setStaticValidators(res http.ResponseWriter, etag string, info fs.FileInfo) {
	res.Header().Set("ETag", etag)

	modified := info.ModTime()
	if !modified.IsZero() {
		res.Header().Set("Last-Modified", modified.UTC().Format(http.TimeFormat))
	}
}

func staticModifiedSinceMatches(req *http.Request, info fs.FileInfo) bool {
	value := req.Header.Get("If-Modified-Since")
	if strings.IsEmpty(value) {
		return false
	}

	since, err := http.ParseTime(value)
	if err != nil {
		return false
	}

	modified := info.ModTime()
	if modified.IsZero() {
		return false
	}

	modified = modified.UTC().Truncate(time.Second.Duration())
	return !modified.After(since)
}

func setStaticContentLength(res http.ResponseWriter, size int64) {
	if size >= 0 {
		res.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
}

func staticETag(name string, info fs.FileInfo) string {
	// W/ marks a weak ETag: metadata freshness, not byte-for-byte content identity.
	return strings.Concat(
		`W/"`,
		name,
		"-",
		strconv.FormatInt(info.Size(), 10),
		"-",
		strconv.FormatInt(info.ModTime().UTC().UnixNano(), 10),
		`"`,
	)
}

func staticETagMatches(value, etag string) bool {
	tag := staticETagOpaqueTag(etag)
	for {
		candidate, rest, found := strings.Cut(value, ",")
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || staticETagOpaqueTag(candidate) == tag {
			return true
		}
		if !found {
			return false
		}

		value = rest
	}
}

func staticETagOpaqueTag(value string) string {
	if strings.HasPrefix(value, "W/") {
		return value[2:]
	}

	return value
}
