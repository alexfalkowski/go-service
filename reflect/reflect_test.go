package reflect_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/reflect"
	"github.com/stretchr/testify/require"
)

func TestIsNil(t *testing.T) {
	require.True(t, reflect.IsNil(nil))
}

func TestIsNilWithTypedNil(t *testing.T) {
	var err error = (*nilError)(nil)

	require.True(t, reflect.IsNil(err))
}

func TestIsNilWithValue(t *testing.T) {
	require.False(t, reflect.IsNil("value"))
}

func TestIsNilWithNonNillableKind(t *testing.T) {
	require.False(t, reflect.IsNil(42))
}

type nilError struct{}

func (e *nilError) Error() string {
	return "nil"
}
