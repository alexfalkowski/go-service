package errors_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestSafeMessage(t *testing.T) {
	for _, tc := range []struct {
		err  error
		name string
		want string
	}{
		{name: "nil", want: "fallback"},
		{name: "plain", err: errors.New("internal"), want: "fallback"},
		{name: "wrapped", err: fmt.Errorf("wrapped: %w", test.SafeMessageError{Message: "safe"}), want: "safe"},
		{name: "empty", err: test.EmptySafeMessageError{Err: test.SafeMessageError{Message: "safe"}}, want: "safe"},
		{name: "joined", err: errors.Join(errors.New("internal"), test.SafeMessageError{Message: "safe"}), want: "safe"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, errors.SafeMessage(tc.err, "fallback"))
		})
	}
}
