package status

import "github.com/alexfalkowski/go-service/v2/context"

const requestErrorKey = context.Key("http-status-request-error")

type requestError struct {
	err error
}

// WithRequestError returns a context that captures the request diagnostic error.
//
// HTTP logging middleware uses the captured error to attach operator diagnostics to the request log while
// [WriteError] continues to render only the safe client message.
func WithRequestError(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestErrorKey, &requestError{})
}

// RequestError returns the captured request diagnostic error for ctx.
//
// It returns nil when ctx was not prepared with [WithRequestError] or no error was written.
func RequestError(ctx context.Context) error {
	state, _ := ctx.Value(requestErrorKey).(*requestError)
	if state == nil {
		return nil
	}

	return state.err
}

// RecordError captures err for request-scoped operator diagnostics without writing a response.
//
// It retains the first error recorded through ctx and does nothing when ctx was not prepared with
// [WithRequestError].
func RecordError(ctx context.Context, err error) {
	if state, _ := ctx.Value(requestErrorKey).(*requestError); state != nil && state.err == nil {
		state.err = err
	}
}
