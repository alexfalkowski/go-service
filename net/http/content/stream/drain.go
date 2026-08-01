package stream

import "github.com/alexfalkowski/go-service/v2/context"

func withDrain(ctx context.Context, drain <-chan struct{}, onDrain func()) (context.Context, context.CancelCauseFunc) {
	if drain == nil {
		return ctx, func(error) {}
	}

	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		select {
		case <-drain:
			cancel(ErrDraining)
			onDrain()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func isDraining(drain <-chan struct{}) bool {
	select {
	case <-drain:
		return true
	default:
		return false
	}
}
