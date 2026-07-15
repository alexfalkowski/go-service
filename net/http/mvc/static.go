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
)

// StaticFile registers an HTTP GET route that serves the named file from the registered filesystem.
//
// It returns false when MVC is not defined (see IsDefined).
func StaticFile(pattern, name string, opts ...StaticOption) bool {
	if !IsDefined() {
		return false
	}

	options := options(opts...)
	handler := func(res http.ResponseWriter, _ *http.Request) {
		serveFile(res, name, options)
	}

	router.Handle(strings.Join(strings.Space, http.MethodGet, pattern), http.HandlerFunc(handler))
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
		serveFile(res, name, options)
	}

	router.Handle(strings.Join(strings.Space, http.MethodGet, pattern), http.HandlerFunc(handler))
	return true
}

func serveFile(res http.ResponseWriter, name string, options *staticOptions) {
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

func setStaticContentLength(res http.ResponseWriter, size int64) {
	if size >= 0 {
		res.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
}
