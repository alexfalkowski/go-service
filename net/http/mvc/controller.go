package mvc

import "github.com/alexfalkowski/go-service/v2/context"

// Controller executes an MVC action and returns a View, a model, and an error.
//
// Controllers are invoked by the routing helpers in this package (see Route/Get/Post/etc.). Those helpers
// populate the request context with HTTP request/response values using `net/http/meta.WithRequest` and
// `net/http/meta.WithResponse` before calling the Controller.
//
// Return values:
//   - view: the view that should be rendered. It is typically created via NewFullView/NewPartialView or NewViewPair.
//   - model: the model passed to the template when rendering succeeds.
//   - err: an error indicating the controller failed.
//
// Error behavior in server wiring:
// When a controller returns a non-nil err, the handler produced by Route writes a status code derived from the
// error (via net/http/status.Code) and renders the returned view using the error value as the template model.
// Implementations should therefore ensure view is non-nil even in error cases, or accept that rendering may panic.
type Controller[Model any] func(ctx context.Context) (*View, *Model, error)
