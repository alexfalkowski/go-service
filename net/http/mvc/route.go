package mvc

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/status"
)

// Controller for mvc.
type Controller[Model any] func(ctx context.Context, req *http.Request, res http.ResponseWriter) (*Model, error)

// Route the template with the path.
func Route[Res any](path string, view *View, controller Controller[Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		r, err := controller(ctx, req, res)
		if err != nil {
			res.WriteHeader(status.Code(err))
			view.ExecuteFailure(res, err)

			return
		}

		err = view.ExecuteSuccess(res, r)
		if err != nil {
			res.WriteHeader(status.Code(err))
			view.ExecuteFailure(res, err)

			return
		}
	}

	mux.HandleFunc(path, h)
}
