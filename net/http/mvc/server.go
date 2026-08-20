package mvc

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
)

var htmlContentType = media.MustParse(media.HTML).WithUTF8()

// NotFound registers controller as the MVC not-found renderer.
//
// It returns false when MVC is not defined (see IsDefined).
func NotFound[Model any](controller NotFoundController[Model]) bool {
	if !IsDefined() {
		return false
	}

	notFoundController = func(ctx context.Context) (*View, any) {
		view, model := controller(ctx)
		return view, model
	}
	return true
}

// Delete registers an HTTP DELETE route that invokes controller.
func Delete[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	return Route(strings.Join(strings.Space, http.MethodDelete, pattern), controller, options...)
}

// Get registers an HTTP GET route that invokes controller.
func Get[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	return Route(strings.Join(strings.Space, http.MethodGet, pattern), controller, options...)
}

// Post registers an HTTP POST route that invokes controller.
func Post[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	return Route(strings.Join(strings.Space, http.MethodPost, pattern), controller, options...)
}

// Put registers an HTTP PUT route that invokes controller.
func Put[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	return Route(strings.Join(strings.Space, http.MethodPut, pattern), controller, options...)
}

// Patch registers an HTTP PATCH route that invokes controller.
func Patch[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	return Route(strings.Join(strings.Space, http.MethodPatch, pattern), controller, options...)
}

// Route registers a handler for pattern that invokes controller and renders the returned view.
//
// It returns false when MVC is not defined (see IsDefined).
// Options apply MVC route registration policy.
//
// The handler sets the response Content-Type to HTML and stores the request and response writer in the
// request context (via net/http/meta) before invoking the controller.
//
// If controller returns an error, the handler renders the returned view using a safe Error model and writes the
// corresponding status code (see [github.com/alexfalkowski/go-service/v2/net/http/status.Code]) only after rendering succeeds. The raw error remains
// available as `mvcModelError` metadata for compatibility; templates that render it can expose diagnostic details.
// If rendering itself fails, the handler writes the render error status instead.
//
// Controller errors and rendering failures are recorded as request-scoped operator diagnostics via
// [github.com/alexfalkowski/go-service/v2/net/http/status.RecordError], surfaced through the HTTP access log.
// When more than one error occurs, the first error is retained.
func Route[Model any](pattern string, controller Controller[Model], options ...RouteOption) bool {
	if !IsDefined() {
		return false
	}

	routeOptions := &routeOptions{}
	for _, option := range options {
		option.apply(routeOptions)
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(http.ContentTypeKey, htmlContentType)

		ctx := req.Context()
		ctx = meta.WithRequestResponse(ctx, req, res)

		view, model, err := controller(ctx)
		if err != nil {
			status.RecordError(ctx, err)

			code := status.Code(err)
			message := errors.SafeMessage(err, status.DefaultMessage(code))
			model := &Error{Code: code, Message: message}

			ctx = meta.WithAttributes(ctx, meta.NewPair("mvcModelError", meta.Error(err)))
			writeView(ctx, res, view, model, model.Code)
			return
		}

		writeView(ctx, res, view, model, http.StatusOK)
	}

	router.HandleRoute(pattern, http.HandlerFunc(handler), routeOptions.httpOptions()...)
	return true
}

func writeNotFound(req *http.Request, res http.ResponseWriter) {
	err := status.Error(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	res.Header().Set(http.ContentTypeKey, htmlContentType)
	ctx := req.Context()
	ctx = meta.WithRequestResponse(ctx, req, res)
	ctx = meta.WithAttributes(ctx, meta.NewPair("mvcModelError", meta.Error(err)))

	view, model := notFoundController(ctx)
	if err := renderView(ctx, res, view, model, http.StatusNotFound); err != nil {
		res.WriteHeader(status.Code(err))
	}
}
