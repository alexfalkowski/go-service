package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/ptr"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/stretchr/testify/require"
)

func TestAppendEmpty(t *testing.T) {
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

	t.Run("typed nil interface", func(t *testing.T) {
		var writer *test.ErrWriter
		var elem io.Writer = writer

		require.Empty(t, slices.AppendNotZero([]io.Writer{}, elem))
	})

	t.Run("nil slice", func(t *testing.T) {
		var elem []string

		require.Empty(t, slices.AppendNotZero([][]string{}, elem))
	})

	t.Run("zero non-comparable struct", func(t *testing.T) {
		require.Empty(t, slices.AppendNotZero([]test.Page{}, test.Page{}))
	})
}

func TestAppendNotZero(t *testing.T) {
	integer := 2

	for _, elem := range []*int{&integer} {
		t.Run("non-zero pointer", func(t *testing.T) {
			require.NotEmpty(t, slices.AppendNotZero([]*int{}, elem))
		})
	}

	for _, elem := range []*int{&integer} {
		t.Run("non-nil pointer", func(t *testing.T) {
			require.NotEmpty(t, slices.AppendNotZero([]*int{}, elem))
		})
	}

	t.Run("empty slice", func(t *testing.T) {
		require.NotEmpty(t, slices.AppendNotZero([][]string{}, []string{}))
	})

	t.Run("non-zero non-comparable struct", func(t *testing.T) {
		value := test.Page{Title: "test"}

		require.NotEmpty(t, slices.AppendNotZero([]test.Page{}, value))
	})
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
