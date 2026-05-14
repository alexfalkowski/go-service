package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/stretchr/testify/require"
)

func TestEmptyAppendZero(t *testing.T) {
	for _, elem := range []*int{nil} {
		t.Run("nil pointer", func(t *testing.T) {
			require.Empty(t, slices.AppendNotZero([]*int{}, elem))
		})
	}

	for _, elem := range []int{0} {
		t.Run("zero value", func(t *testing.T) {
			require.Empty(t, slices.AppendNotZero([]int{}, elem))
		})
	}
}

func TestEmptyAppendNil(t *testing.T) {
	for _, elem := range []*int{nil} {
		t.Run("nil pointer", func(t *testing.T) {
			require.Empty(t, slices.AppendNotNil([]*int{}, elem))
		})
	}
}

func TestAppendZero(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		t.Run("non-zero pointer", func(t *testing.T) {
			require.NotEmpty(t, slices.AppendNotZero([]*int{}, elem))
		})
	}
}

func TestAppendNil(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		t.Run("non-nil pointer", func(t *testing.T) {
			require.NotEmpty(t, slices.AppendNotZero([]*int{}, elem))
		})
	}
}

func TestElemFunc(t *testing.T) {
	elems := []*string{ptr.Value("test")}

	elem, ok := slices.ElemFunc(elems, func(t *string) bool { return *t == "test" })
	require.NotNil(t, elem)
	require.True(t, ok)

	elem, ok = slices.ElemFunc(elems, func(t *string) bool { return *t == "bob" })
	require.Nil(t, elem)
	require.False(t, ok)
}

func TestClip(t *testing.T) {
	values := [...]string{"first", "second"}
	first := slices.Clip(values[0:1])
	second := values[1:2]

	first = append(first, "next")

	require.Equal(t, []string{"first", "next"}, first)
	require.Equal(t, []string{"second"}, second)
}
