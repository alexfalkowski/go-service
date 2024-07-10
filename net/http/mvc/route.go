package mvc

import (
	"context"
	"net/http"
)

// Controller for mvc.
type Controller[Model any] func(ctx context.Context, req *http.Request, res http.ResponseWriter) (*Model, error)

// Route the template with the path.
func Route[Res any](path string, view *View, controller Controller[Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var response *Res

		if controller != nil {
			r, err := controller(ctx, req, res)
			if err != nil {
				WriteError(ctx, view, res, err)

				return
			}

			response = r
		}

		if err := view.ExecuteSuccess(res, response); err != nil {
			WriteError(ctx, view, res, err)

			return
		}
	}

	mux.HandleFunc(path, h)
}
