package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

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
