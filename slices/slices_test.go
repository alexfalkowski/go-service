package slices_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/stretchr/testify/require"
)

func TestAppendEmpty(t *testing.T) {
	t.Parallel()

	var writer *test.ErrWriter
	var elem io.Writer = writer
	var nilSlice []string

	tests := []struct {
		append func() any
		name   string
	}{
		{name: "nil pointer", append: func() any { return slices.AppendNotZero([]*int{}, nil) }},
		{name: "zero value", append: func() any { return slices.AppendNotZero([]int{}, 0) }},
		{name: "typed nil interface", append: func() any { return slices.AppendNotZero([]io.Writer{}, elem) }},
		{name: "nil slice", append: func() any { return slices.AppendNotZero([][]string{}, nilSlice) }},
		{name: "zero non-comparable struct", append: func() any { return slices.AppendNotZero([]test.Page{}, test.Page{}) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Empty(t, test.append())
		})
	}
}

func TestAppendNotZero(t *testing.T) {
	t.Parallel()

	integer := 2

	tests := []struct {
		append func() any
		name   string
	}{
		{name: "non-zero pointer", append: func() any { return slices.AppendNotZero([]*int{}, &integer) }},
		{name: "empty slice", append: func() any { return slices.AppendNotZero([][]string{}, []string{}) }},
		{name: "non-zero non-comparable struct", append: func() any {
			return slices.AppendNotZero([]test.Page{}, test.Page{Title: "test"})
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.NotEmpty(t, test.append())
		})
	}
}

func TestBackward(t *testing.T) {
	t.Parallel()

	values := []string{}
	for _, value := range slices.Backward([]string{"first", "second"}) {
		values = append(values, value)
	}

	require.Equal(t, []string{"second", "first"}, values)
}

func TestElemFunc(t *testing.T) {
	t.Parallel()

	elems := []*string{new("test")}

	tests := []struct {
		name  string
		value string
		found bool
	}{
		{name: "match", value: "test", found: true},
		{name: "missing", value: "bob"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			elem, ok := slices.ElemFunc(elems, func(value *string) bool { return *value == test.value })

			require.Equal(t, test.found, ok)
			if !test.found {
				require.Nil(t, elem)
				return
			}

			require.NotNil(t, elem)
			require.Equal(t, test.value, *elem)
		})
	}
}

func TestClip(t *testing.T) {
	t.Parallel()

	values := [...]string{"first", "second"}
	first := slices.Clip(values[0:1])
	second := values[1:2]

	first = append(first, "next")

	require.Equal(t, []string{"first", "next"}, first)
	require.Equal(t, []string{"second"}, second)
}

func TestContains(t *testing.T) {
	t.Parallel()

	type names []string

	t.Run("present", func(t *testing.T) {
		t.Parallel()

		require.True(t, slices.Contains(names{"alice", "bob"}, "bob"))
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		require.False(t, slices.Contains(names{"alice", "bob"}, "eve"))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		require.False(t, slices.Contains(names{}, "alice"))
	})
}
