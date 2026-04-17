package context_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestBackground(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	deadline, ok := ctx.Deadline()
	require.False(t, ok)
	require.True(t, deadline.IsZero())
	require.Nil(t, ctx.Done())
	require.NoError(t, ctx.Err())
	require.NoError(t, context.Cause(ctx))
}

func TestCauseBeforeCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, context.Cause(ctx))
}

func TestWithCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	<-ctx.Done()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
	require.ErrorIs(t, context.Cause(ctx), context.Canceled)
}

func TestWithCancelCause(t *testing.T) {
	t.Parallel()

	t.Run("records provided cause", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("cancel cause")
		ctx, cancel := context.WithCancelCause(context.Background())

		cancel(cause)

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.ErrorIs(t, context.Cause(ctx), cause)
	})

	t.Run("falls back to canceled when nil cause", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancelCause(context.Background())

		cancel(nil)

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.ErrorIs(t, context.Cause(ctx), context.Canceled)
	})
}

func TestWithTimeoutCause(t *testing.T) {
	t.Parallel()

	t.Run("returns configured cause when timeout expires", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("timeout cause")
		ctx, cancel := context.WithTimeoutCause(context.Background(), 0, cause)
		defer cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		require.ErrorIs(t, context.Cause(ctx), cause)
	})

	t.Run("manual cancel keeps canceled as cause", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("timeout cause")
		ctx, cancel := context.WithTimeoutCause(context.Background(), time.Minute, cause)

		cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.ErrorIs(t, context.Cause(ctx), context.Canceled)
	})
}

func TestWithDeadline(t *testing.T) {
	t.Parallel()

	t.Run("returns deadline exceeded when deadline passes", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Millisecond.Duration()))
		defer cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		require.ErrorIs(t, context.Cause(ctx), context.DeadlineExceeded)
	})

	t.Run("manual cancel keeps canceled as cause", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute.Duration()))

		cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.ErrorIs(t, context.Cause(ctx), context.Canceled)
	})
}

func TestWithDeadlineCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("deadline cause")
	ctx, cancel := context.WithDeadlineCause(context.Background(), time.Now().Add(-time.Millisecond.Duration()), cause)
	defer cancel()

	<-ctx.Done()
	require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	require.ErrorIs(t, context.Cause(ctx), cause)
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	t.Run("returns deadline exceeded when timeout expires", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		require.ErrorIs(t, context.Cause(ctx), context.DeadlineExceeded)
	})

	t.Run("manual cancel keeps canceled as cause", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

		cancel()

		<-ctx.Done()
		require.ErrorIs(t, ctx.Err(), context.Canceled)
		require.ErrorIs(t, context.Cause(ctx), context.Canceled)
	})
}

func TestCausePropagatesFromParent(t *testing.T) {
	t.Parallel()

	cause := errors.New("parent cause")
	parent, cancel := context.WithCancelCause(context.Background())
	child, childCancel := context.WithTimeout(parent, time.Minute)
	defer childCancel()

	cancel(cause)

	<-child.Done()
	require.ErrorIs(t, parent.Err(), context.Canceled)
	require.ErrorIs(t, context.Cause(parent), cause)
	require.ErrorIs(t, context.Cause(child), cause)
}

func TestWithoutCancel(t *testing.T) {
	t.Parallel()

	parent, cancel := context.WithCancelCause(context.WithValue(context.Background(), context.Key("key"), "value"))
	child := context.WithoutCancel(parent)

	cancel(errors.New("parent cause"))

	deadline, ok := child.Deadline()
	require.False(t, ok)
	require.True(t, deadline.IsZero())
	require.Nil(t, child.Done())
	require.NoError(t, child.Err())
	require.NoError(t, context.Cause(child))
	require.Equal(t, "value", child.Value(context.Key("key")))
}

func TestWithValue(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), context.Key("key"), "value")

	require.Equal(t, "value", ctx.Value(context.Key("key")))
	require.Nil(t, ctx.Value(context.Key("missing")))
}
