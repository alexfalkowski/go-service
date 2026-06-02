package reflect_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/reflect"
	"github.com/stretchr/testify/require"
)

func TestIsNil(t *testing.T) {
	var err error = (*test.NilError)(nil)

	tests := []struct {
		value any
		name  string
		want  bool
	}{
		{name: "nil", value: nil, want: true},
		{name: "typed nil", value: err, want: true},
		{name: "value", value: "value"},
		{name: "kind that cannot be nil", value: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, reflect.IsNil(tt.value))
		})
	}
}

func TestIsZero(t *testing.T) {
	var err error = (*test.NilError)(nil)

	tests := []struct {
		value any
		name  string
		want  bool
	}{
		{name: "nil", value: nil, want: true},
		{name: "typed nil", value: err, want: true},
		{name: "value", value: "value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, reflect.IsZero(tt.value))
		})
	}
}
