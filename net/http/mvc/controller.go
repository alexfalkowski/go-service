package mvc

import "context"

// Controller for mvc.
type Controller[Model any] func(ctx context.Context) (View, *Model, error)
