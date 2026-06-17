package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

// testBytesSink keeps conversion results observable for allocation assertions.
var testBytesSink []byte

func TestBytes(t *testing.T) {
	tests := []struct {
		name string
		s    string
	}{
		{name: "empty", s: ""},
		{name: "non empty", s: "hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := strings.Bytes(test.s)

			require.Len(t, got, len(test.s))
			if len(test.s) == 0 {
				require.Empty(t, got)
			} else {
				require.Equal(t, []byte(test.s), got)
			}
		})
	}
}

func TestBytesDoesNotAllocate(t *testing.T) {
	s := "hello"
	allocs := testing.AllocsPerRun(100, func() {
		testBytesSink = strings.Bytes(s)
	})

	require.Zero(t, allocs)
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{name: "empty", s: "", want: true},
		{name: "space", s: " ", want: false},
		{name: "populated", s: "request", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.IsEmpty(test.s))
		})
	}
}

func TestIsAnyEmpty(t *testing.T) {
	tests := []struct {
		name string
		ss   []string
		want bool
	}{
		{name: "no strings", ss: []string{}, want: false},
		{name: "all populated", ss: []string{"request", "response"}, want: false},
		{name: "one empty", ss: []string{"request", "", "response"}, want: true},
		{name: "all empty", ss: []string{"", ""}, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.IsAnyEmpty(test.ss...))
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		sep  string
		want string
		ss   []string
	}{
		{name: "no strings", sep: "/", ss: []string{}, want: ""},
		{name: "one string", sep: "/", ss: []string{"service"}, want: "service"},
		{name: "multiple strings", sep: "/", ss: []string{"service", "admin", "metrics"}, want: "service/admin/metrics"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.Join(test.sep, test.ss...))
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name string
		want string
		ss   []string
	}{
		{name: "no strings", ss: []string{}, want: ""},
		{name: "one string", ss: []string{"service"}, want: "service"},
		{name: "multiple strings", ss: []string{"service", "admin", "metrics"}, want: "serviceadminmetrics"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.Concat(test.ss...))
		})
	}
}

func TestCutColon(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantBefore string
		wantAfter  string
		wantFound  bool
	}{
		{name: "env source", s: "env:NAME", wantBefore: "env", wantAfter: "NAME", wantFound: true},
		{name: "file source with later colon", s: "file:/tmp/a:b", wantBefore: "file", wantAfter: "/tmp/a:b", wantFound: true},
		{name: "missing colon", s: "literal", wantBefore: "literal", wantAfter: ""},
		{name: "leading colon", s: ":value", wantBefore: "", wantAfter: "value", wantFound: true},
		{name: "trailing colon", s: "env:", wantBefore: "env", wantAfter: "", wantFound: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before, after, found := strings.CutColon(test.s)

			require.Equal(t, test.wantBefore, before)
			require.Equal(t, test.wantAfter, after)
			require.Equal(t, test.wantFound, found)
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{name: "last match", s: "/service/admin/metrics", substr: "/", want: 14},
		{name: "missing", s: "service", substr: "/", want: -1},
		{name: "empty substr", s: "service", substr: "", want: 7},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.LastIndex(test.s, test.substr))
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{name: "matches", s: "/greet.v1.Greeter/SayHello", substr: "/", want: 2},
		{name: "missing", s: "greet.v1.Greeter", substr: "/", want: 0},
		{name: "empty substr", s: "service", substr: "", want: 8},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.Count(test.s, test.substr))
		})
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		name string
		s    string
		old  string
		new  string
		want string
	}{
		{name: "replaces all matches", s: "service/admin/service", old: "service", new: "api", want: "api/admin/api"},
		{name: "missing match", s: "service", old: "admin", new: "api", want: "service"},
		{name: "empty old string", s: "api", old: "", new: "/", want: "/a/p/i/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.ReplaceAll(test.s, test.old, test.new))
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		cutset string
		want   string
	}{
		{name: "trims cutset", s: "/service/admin/", cutset: "/", want: "service/admin"},
		{name: "trims repeated cutset", s: "::service::", cutset: ":", want: "service"},
		{name: "missing cutset", s: "service", cutset: "/", want: "service"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, strings.Trim(test.s, test.cutset))
		})
	}
}
