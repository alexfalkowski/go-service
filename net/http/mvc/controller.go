package mvc

import "github.com/alexfalkowski/go-service/v2/context"

// Controller executes an MVC action and returns a View, a model, and an error.
type Controller[Model any] func(ctx context.Context) (*View, *Model, error)
