package errors_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestNewHandlerNilLogger(t *testing.T) {
	require.Nil(t, errors.NewHandler(nil))
}

func TestRegisterNilHandler(t *testing.T) {
	original := otel.GetErrorHandler()
	defer otel.SetErrorHandler(original)

	errors.Register(nil)

	require.Same(t, original, otel.GetErrorHandler())
}

func TestHandleNilHandler(t *testing.T) {
	var handler *errors.Handler

	require.NotPanics(t, func() {
		handler.Handle(context.Canceled)
	})
}
