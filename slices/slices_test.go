package slices_test

import (
	"cmp"
	"testing"

	"github.com/alexfalkowski/go-service/v3/internal/test"
	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/alexfalkowski/go-service/v3/slices"
	"github.com/stretchr/testify/require"
)

func TestAppendNotZeroExcludesZeroValues(t *testing.T) {
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

func TestAppendNotZeroIncludesNonZeroValues(t *testing.T) {
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

func TestBackwardIteratesValuesInReverseOrder(t *testing.T) {
	t.Parallel()

	values := []string{}
	for _, value := range slices.Backward([]string{"first", "second"}) {
		values = append(values, value)
	}

	require.Equal(t, []string{"second", "first"}, values)
}

func TestElemFuncReturnsMatchingElementOrNil(t *testing.T) {
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

func TestClipPreventsAppendFromMutatingOriginalArray(t *testing.T) {
	t.Parallel()

	values := [...]string{"first", "second"}
	first := slices.Clip(values[0:1])
	second := values[1:2]

	first = append(first, "next")

	require.Equal(t, []string{"first", "next"}, first)
	require.Equal(t, []string{"second"}, second)
}

func TestContainsRecognizesPresentValue(t *testing.T) {
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

func TestSortFuncUsesCustomOrder(t *testing.T) {
	t.Parallel()

	values := []int{2, 1, 3}
	slices.SortFunc(values, func(a, b int) int { return cmp.Compare(b, a) })

	require.Equal(t, []int{3, 2, 1}, values)
	require.True(t, slices.IsSortedFunc(values, func(a, b int) int { return cmp.Compare(b, a) }), "sorted by reverse order")
	require.False(t, slices.IsSortedFunc(values, cmp.Compare[int]), "unsorted by natural order")
	require.False(t, slices.IsSorted(values))
}

func TestStableSortPreservesEqualElementOrder(t *testing.T) {
	t.Parallel()

	compare := func(a, b string) int { return cmp.Compare(len(a), len(b)) }
	tests := []struct {
		sort func([]string) []string
		name string
	}{
		{name: "slice", sort: func(values []string) []string {
			slices.SortStableFunc(values, compare)
			return values
		}},
		{name: "iterator", sort: func(values []string) []string {
			return slices.SortedStableFunc(slices.Values(values), compare)
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			values := []string{"aa", "c", "bb", "d", "cc", "e", "dd", "f", "ee", "g", "ff", "h", "gg", "i", "hh", "j", "ii", "k", "jj", "l"}
			want := []string{"c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj"}

			require.Equal(t, want, test.sort(values))
		})
	}
}

func TestSortedCollectsAndOrdersIterator(t *testing.T) {
	t.Parallel()

	values := []int{2, 1, 3}
	ascending := slices.Sorted(slices.Values(values))
	descending := slices.SortedFunc(slices.Values(values), func(a, b int) int { return cmp.Compare(b, a) })

	require.Equal(t, []int{1, 2, 3}, ascending)
	require.True(t, slices.IsSorted(ascending))
	require.Equal(t, []int{3, 2, 1}, descending)
	require.Equal(t, []int{2, 1, 3}, values, "sorting an iterator preserves its source slice")
}

func TestAllIteratesIndicesAndValues(t *testing.T) {
	t.Parallel()

	indices := []int{}
	values := []string{}
	for index, value := range slices.All([]string{"first", "second"}) {
		indices = append(indices, index)
		values = append(values, value)
	}

	require.Equal(t, []int{0, 1}, indices)
	require.Equal(t, []string{"first", "second"}, values)
}

func TestAppendSeqAndChunkPreserveNamedSliceType(t *testing.T) {
	t.Parallel()

	type numbers []int

	values := slices.AppendSeq(numbers{1}, slices.Values([]int{2, 3, 4, 5}))
	chunks := slices.Collect(slices.Chunk(values, 2))

	require.Equal(t, numbers{1, 2, 3, 4, 5}, values)
	require.Equal(t, []numbers{{1, 2}, {3, 4}, {5}}, chunks)
}

func TestSliceChangesPreserveNamedSliceType(t *testing.T) {
	t.Parallel()

	type numbers []int

	tests := []struct {
		change func() any
		name   string
		want   numbers
	}{
		{name: "compact", change: func() any { return slices.Compact(numbers{1, 1, 2, 2, 1}) }, want: numbers{1, 2, 1}},
		{name: "compact func", change: func() any {
			return slices.CompactFunc(numbers{1, -1, 2, -2}, func(a, b int) bool { return a == b || a == -b })
		}, want: numbers{1, 2}},
		{name: "concat", change: func() any { return slices.Concat(numbers{1}, numbers{2, 3}) }, want: numbers{1, 2, 3}},
		{name: "delete", change: func() any { return slices.Delete(numbers{1, 2, 3, 4}, 1, 3) }, want: numbers{1, 4}},
		{name: "insert", change: func() any { return slices.Insert(numbers{1, 4}, 1, 2, 3) }, want: numbers{1, 2, 3, 4}},
		{name: "repeat", change: func() any { return slices.Repeat(numbers{1, 2}, 2) }, want: numbers{1, 2, 1, 2}},
		{name: "replace", change: func() any { return slices.Replace(numbers{1, 2, 3, 4}, 1, 3, 5) }, want: numbers{1, 5, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, test.change())
		})
	}
}

func TestGrowPreservesValuesAndReservesCapacity(t *testing.T) {
	t.Parallel()

	type numbers []int

	values := slices.Grow(numbers{1, 2}, 10)

	require.Equal(t, numbers{1, 2}, values)
	require.GreaterOrEqual(t, cap(values)-len(values), 10)
}

func TestReverseChangesSliceInPlace(t *testing.T) {
	t.Parallel()

	values := []int{1, 2, 3}
	slices.Reverse(values)

	require.Equal(t, []int{3, 2, 1}, values)
}

func TestCompareAndEqualUseElementValues(t *testing.T) {
	t.Parallel()

	values := []int{1, 2}

	require.Negative(t, slices.Compare(values, []int{1, 3}))
	require.Positive(t, slices.Compare(values, []int{1, 1}))
	require.Zero(t, slices.Compare(values, []int{1, 2}))
	require.True(t, slices.Equal(values, []int{1, 2}), "identical elements are equal")
	require.False(t, slices.Equal(values, []int{2, 1}), "element order matters")

	compare := func(number int, word string) int { return cmp.Compare(number, len(word)) }
	equal := func(number int, word string) bool { return number == len(word) }

	require.Negative(t, slices.CompareFunc([]int{1, 2}, []string{"a", "bbb"}, compare))
	require.Positive(t, slices.CompareFunc([]int{1, 2}, []string{"a", "b"}, compare))
	require.Zero(t, slices.CompareFunc([]int{1, 2}, []string{"a", "bb"}, compare))
	require.True(t, slices.EqualFunc([]int{1, 2}, []string{"a", "bb"}, equal), "matching lengths are equal")
	require.False(t, slices.EqualFunc([]int{1, 2}, []string{"a", "bbb"}, equal), "different lengths are unequal")
}

func TestIndexReportsFirstMatchOrMissing(t *testing.T) {
	t.Parallel()

	values := []int{1, 2, 2}
	tests := []struct {
		name   string
		target int
		want   int
	}{
		{name: "first match", target: 2, want: 1},
		{name: "missing", target: 3, want: -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, slices.Index(values, test.target), "value search")
			require.Equal(t, test.want, slices.IndexFunc(values, func(value int) bool { return value == test.target }), "predicate search")
		})
	}
}

func TestBinarySearchReportsMatchAndInsertionPoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target int
		index  int
		found  bool
	}{
		{name: "match", target: 3, index: 1, found: true},
		{name: "insert between", target: 2, index: 1},
		{name: "insert after", target: 6, index: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			index, found := slices.BinarySearch([]int{1, 3, 5}, test.target)
			require.Equal(t, test.index, index, "ordered search index")
			require.Equal(t, test.found, found, "ordered search match")

			index, found = slices.BinarySearchFunc([]string{"a", "bbb", "ccccc"}, test.target, func(value string, target int) int {
				return cmp.Compare(len(value), target)
			})
			require.Equal(t, test.index, index, "custom search index")
			require.Equal(t, test.found, found, "custom search match")
		})
	}
}

func TestExtremaHonorCustomComparison(t *testing.T) {
	t.Parallel()

	values := []int{2, 1, 3}
	compare := func(a, b int) int { return cmp.Compare(b, a) }

	require.Equal(t, 1, slices.Min(values), "minimum by natural order")
	require.Equal(t, 3, slices.Max(values), "maximum by natural order")
	require.Equal(t, 3, slices.MinFunc(values, compare), "minimum by reverse order")
	require.Equal(t, 1, slices.MaxFunc(values, compare), "maximum by reverse order")
}
