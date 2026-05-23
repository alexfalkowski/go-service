package reflect_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/reflect"
	"github.com/stretchr/testify/require"
)

func TestIsNil(t *testing.T) {
	require.True(t, reflect.IsNil(nil))
}

func TestIsNilWithTypedNil(t *testing.T) {
	var err error = (*test.NilError)(nil)

	require.True(t, reflect.IsNil(err))
}

func TestIsNilWithValue(t *testing.T) {
	require.False(t, reflect.IsNil("value"))
}

func TestIsNilWithNonNillableKind(t *testing.T) {
	require.False(t, reflect.IsNil(42))
}

func TestIsZero(t *testing.T) {
	require.True(t, reflect.IsZero(nil))
}

func TestIsZeroWithTypedNil(t *testing.T) {
	var err error = (*test.NilError)(nil)

	require.True(t, reflect.IsZero(err))
}

func TestIsZeroWithValue(t *testing.T) {
	require.False(t, reflect.IsZero("value"))
}
