package mvc

import "github.com/alexfalkowski/go-service/v2/context"

// Controller for mvc.
type Controller[Model any] func(ctx context.Context) (*View, *Model, error)
