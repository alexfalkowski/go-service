package mvc

import "github.com/alexfalkowski/go-service/v2/context"

// Controller executes an MVC action and returns a View, a model, and an error.
//
// Controllers are invoked by the routing helpers in this package (see Route/Get/Post/etc.). Those helpers
// populate the request context with HTTP content metadata before calling the Controller.
//
// Return values:
//   - view: the view that should be rendered. It is typically prebuilt during startup or route registration
//     via NewFullView/NewPartialView or NewViewPair.
//   - model: the model passed to the template when rendering succeeds.
//   - err: an error indicating the controller failed.
//
// Error behavior in server wiring:
// When a controller returns a non-nil err, the handler produced by Route writes a status code derived from the
// error (via net/http/status.Code) and renders the returned view using a client-safe Error model.
// The `mvcModelError` metadata value contains the raw error string for compatibility; templates that render it
// can expose diagnostic details.
// Implementations should therefore return a non-nil view when they want to render an error page. If view is nil,
// rendering returns ErrMissingView and the route writes the status derived from that error instead.
type Controller[Model any] func(ctx context.Context) (*View, *Model, error)

// NotFoundController returns the view and model used to render an MVC not-found response.
//
// It is used by Handler for unmatched routes.
type NotFoundController[Model any] func(ctx context.Context) (*View, *Model)
