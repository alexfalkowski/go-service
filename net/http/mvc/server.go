package mvc

import (
	"io"
	"io/fs"
	"path"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete registers an HTTP DELETE route that invokes controller.
func Delete[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodDelete, pattern), controller)
}

// Get registers an HTTP GET route that invokes controller.
func Get[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodGet, pattern), controller)
}

// Post registers an HTTP POST route that invokes controller.
func Post[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPost, pattern), controller)
}

// Put registers an HTTP PUT route that invokes controller.
func Put[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPut, pattern), controller)
}

// Patch registers an HTTP PATCH route that invokes controller.
func Patch[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPatch, pattern), controller)
}

// Route registers a handler for pattern that invokes controller and renders the returned view.
//
// It returns false when MVC is not defined (see IsDefined).
//
// The handler sets the response Content-Type to HTML and stores the request and response writer in the
// request context (via net/http/meta) before invoking the controller.
//
// If controller returns an error, the handler renders the returned view using the error value as the template
// model and writes the corresponding status code (see net/http/status.Code) only after rendering succeeds.
// If rendering itself fails, the handler writes the render error status instead.
func Route[Model any](pattern string, controller Controller[Model]) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(content.TypeKey, media.HTML)

		ctx := req.Context()
		ctx = meta.WithContent(ctx, req, res, nil)

		view, model, err := controller(ctx)
		if err != nil {
			ctx = meta.WithAttributes(ctx, meta.NewPair("mvcModelError", meta.Error(err)))
			writeView(ctx, res, view, err, status.Code(err))
			return
		}

		writeView(ctx, res, view, model, http.StatusOK)
	}

	http.HandleFunc(mux, pattern, handler)
	return true
}

// StaticFile registers an HTTP GET route that serves the named file from the registered filesystem.
//
// It returns false when MVC is not defined (see IsDefined).
func StaticFile(pattern, name string) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		buffer := pool.Get()
		defer pool.Put(buffer)

		ctx := req.Context()

		if err := writeFile(name, buffer); err != nil {
			writeHeader(ctx, res, staticStatusCode(err))
			return
		}

		writeBuffer(ctx, res, http.StatusOK, buffer)
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
func StaticPathValue(pattern, value, prefix string) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		buffer := pool.Get()
		defer pool.Put(buffer)

		ctx := req.Context()
		cleaned := path.Clean(req.PathValue(value))
		if cleaned == "." || cleaned != req.PathValue(value) || !fs.ValidPath(cleaned) || strings.Contains(cleaned, `\`) {
			writeHeader(ctx, res, staticStatusCode(status.BadRequestError(fs.ErrInvalid)))
			return
		}

		name := path.Join(prefix, cleaned)
		if err := writeFile(name, buffer); err != nil {
			writeHeader(ctx, res, staticStatusCode(err))
			return
		}

		writeBuffer(ctx, res, http.StatusOK, buffer)
	}

	http.HandleFunc(mux, strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

func writeFile(name string, writer io.Writer) error {
	f, err := fileSystem.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	return err
}

func staticStatusCode(err error) int {
	if errors.Is(err, fs.ErrNotExist) {
		return http.StatusNotFound
	}
	return status.Code(err)
}

func writeView(ctx context.Context, res http.ResponseWriter, view *View, model any, code int) {
	buffer := pool.Get()
	defer pool.Put(buffer)

	if err := view.render(ctx, buffer, model); err != nil {
		writeHeader(ctx, res, status.Code(err))
		return
	}

	writeBuffer(ctx, res, code, buffer)
}

func writeBuffer(ctx context.Context, res http.ResponseWriter, code int, buffer *bytes.Buffer) {
	if !writeHeader(ctx, res, code) {
		return
	}

	_, _ = buffer.WriteTo(res)
}

func writeHeader(ctx context.Context, res http.ResponseWriter, code int) bool {
	if ctx.Err() != nil {
		return false
	}

	res.WriteHeader(code)
	return true
}
