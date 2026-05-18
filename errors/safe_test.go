package errors_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
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
		{name: "wrapped", err: fmt.Errorf("wrapped: %w", safeError{msg: "safe"}), want: "safe"},
		{name: "empty", err: emptySafeError{err: safeError{msg: "safe"}}, want: "safe"},
		{name: "joined", err: errors.Join(errors.New("internal"), safeError{msg: "safe"}), want: "safe"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, errors.SafeMessage(tc.err, "fallback"))
		})
	}
}

type safeError struct {
	msg string
}

func (s safeError) Error() string {
	return "internal"
}

func (s safeError) SafeMessage() string {
	return s.msg
}

type emptySafeError struct {
	err error
}

func (e emptySafeError) Error() string {
	return "internal"
}

func (e emptySafeError) SafeMessage() string {
	return ""
}

func (e emptySafeError) Unwrap() error {
	return e.err
}
