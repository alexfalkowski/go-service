package mvc

import (
	"io/fs"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
	"github.com/go-sprout/sprout"
)

var htmlContentType = media.MustParse(media.HTML).WithUTF8()

// ServerParams defines dependencies used to construct a [Server].
//
// This package uses a single Server value to avoid threading commonly shared dependencies through
// every routing helper call. These dependencies are typically provided by DI wiring.
type ServerParams struct {
	di.In

	// Router registers MVC routes (Route/Get/Post/etc.).
	Router *http.Router

	// FunctionMap is the template function map used when parsing templates.
	//
	// It is typically constructed via go-sprout and used to provide common template helpers.
	FunctionMap sprout.FunctionMap

	// FileSystem is the filesystem used to load template files and static files.
	//
	// It is optional: when nil, MVC is considered not defined and routing helpers return false.
	FileSystem fs.FS `optional:"true"`

	// Pool is the shared buffer pool used to reduce allocations while rendering views and buffering static files.
	Pool *sync.BufferPool

	// Layout defines the base templates used to render full and partial views.
	//
	// It is optional: when nil or invalid (see Layout.IsValid), MVC is considered not defined and
	// routing helpers return false.
	Layout *Layout `optional:"true"`
}

// NewServer constructs a Server from params.
//
// NewServer is expected to be called during application startup (typically via dependency injection).
//
// Definition rules:
// MVC routing/rendering is considered "defined" only when both:
//   - FileSystem is non-nil, and
//   - Layout is valid (see [Layout.IsValid]).
//
// Pool must be non-nil. Standard go-service module wiring provides it from [sync.Module].
//
// If MVC is not defined, routing helpers (Route/Get/Post/etc. and static helpers) return false and do not
// register handlers.
func NewServer(params ServerParams) *Server {
	return &Server{
		router:     params.Router,
		fmap:       params.FunctionMap,
		fileSystem: params.FileSystem,
		pool:       params.Pool,
		layout:     params.Layout,
	}
}

// Server registers MVC routes and renders views for an HTTP router.
type Server struct {
	router             *http.Router
	fmap               sprout.FunctionMap
	fileSystem         fs.FS
	pool               *sync.BufferPool
	layout             *Layout
	notFoundController func(ctx context.Context) (*View, any)
}

// IsDefined reports whether MVC routing and rendering has been configured.
//
// MVC is considered defined only when a FileSystem is available and Layout is valid (non-nil with both
// layout template names set).
func (s *Server) IsDefined() bool {
	return s.fileSystem != nil && s.layout.IsValid()
}

// NotFound registers controller as the MVC not-found renderer.
//
// It returns false when MVC is not defined (see IsDefined).
func (s *Server) NotFound[Model any](controller NotFoundController[Model]) bool {
	if !s.IsDefined() {
		return false
	}

	s.notFoundController = func(ctx context.Context) (*View, any) {
		view, model := controller(ctx)
		return view, model
	}
	return true
}

// Delete registers an HTTP DELETE route that invokes controller.
func (s *Server) Delete[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	return s.Route(strings.Join(strings.Space, http.MethodDelete, pattern), controller, opts...)
}

// Get registers an HTTP GET route that invokes controller.
func (s *Server) Get[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	return s.Route(strings.Join(strings.Space, http.MethodGet, pattern), controller, opts...)
}

// Post registers an HTTP POST route that invokes controller.
func (s *Server) Post[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	return s.Route(strings.Join(strings.Space, http.MethodPost, pattern), controller, opts...)
}

// Put registers an HTTP PUT route that invokes controller.
func (s *Server) Put[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	return s.Route(strings.Join(strings.Space, http.MethodPut, pattern), controller, opts...)
}

// Patch registers an HTTP PATCH route that invokes controller.
func (s *Server) Patch[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	return s.Route(strings.Join(strings.Space, http.MethodPatch, pattern), controller, opts...)
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
func (s *Server) Route[Model any](pattern string, controller Controller[Model], opts ...RouteOption) bool {
	if !s.IsDefined() {
		return false
	}

	options := &routeOptions{}
	for _, opt := range opts {
		opt.apply(options)
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
			s.writeView(ctx, res, view, model, model.Code)
			return
		}

		s.writeView(ctx, res, view, model, http.StatusOK)
	}

	s.router.HandleRoute(pattern, http.HandlerFunc(handler), options.httpOptions()...)
	return true
}

func (s *Server) writeNotFound(req *http.Request, res http.ResponseWriter) {
	err := status.Error(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	res.Header().Set(http.ContentTypeKey, htmlContentType)
	ctx := req.Context()
	ctx = meta.WithRequestResponse(ctx, req, res)
	ctx = meta.WithAttributes(ctx, meta.NewPair("mvcModelError", meta.Error(err)))

	view, model := s.notFoundController(ctx)
	if err := s.renderView(ctx, res, view, model, http.StatusNotFound); err != nil {
		res.WriteHeader(status.Code(err))
	}
}
