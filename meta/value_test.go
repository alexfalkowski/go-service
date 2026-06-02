package meta_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name       string
		wantValue  string
		wantString string
		value      meta.Value
		wantEmpty  bool
	}{
		{name: "string", value: meta.String("visible"), wantValue: "visible", wantString: "visible"},
		{name: "empty string", value: meta.String(""), wantValue: "", wantString: "", wantEmpty: true},
		{name: "blank", value: meta.Blank(), wantValue: "", wantString: "", wantEmpty: true},
		{name: "ignored", value: meta.Ignored("secret"), wantValue: "secret", wantString: ""},
		{name: "redacted", value: meta.Redacted("secret"), wantValue: "secret", wantString: "******"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.wantValue, test.value.Value())
			require.Equal(t, test.wantString, test.value.String())
			require.Equal(t, test.wantEmpty, test.value.IsEmpty())
		})
	}
}

func TestErrorWithNil(t *testing.T) {
	require.Equal(t, meta.Blank(), meta.Error(nil))
}

func TestError(t *testing.T) {
	value := meta.Error(&test.MessageError{Message: "boom"})

	require.Equal(t, "boom", value.Value())
	require.Equal(t, "boom", value.String())
	require.False(t, value.IsEmpty())
}

func TestErrorWithTypedNil(t *testing.T) {
	var err error = (*test.MessageError)(nil)

	require.Equal(t, meta.Blank(), meta.Error(err))
}

func TestStringerConversions(t *testing.T) {
	tests := []struct {
		name       string
		wantValue  string
		wantString string
		value      meta.Value
	}{
		{name: "to string", value: meta.ToString(&test.Stringer{Value: "visible"}), wantValue: "visible", wantString: "visible"},
		{name: "to redacted", value: meta.ToRedacted(&test.Stringer{Value: "secret"}), wantValue: "secret", wantString: "******"},
		{name: "to ignored", value: meta.ToIgnored(&test.Stringer{Value: "secret"}), wantValue: "secret", wantString: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.wantValue, test.value.Value())
			require.Equal(t, test.wantString, test.value.String())
			require.False(t, test.value.IsEmpty())
		})
	}
}

func TestToStringWithNil(t *testing.T) {
	require.Equal(t, meta.Blank(), meta.ToString(nil))
}

func TestToStringWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToString(stringer))
}

func TestToRedactedWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToRedacted(stringer))
}

func TestToIgnoredWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToIgnored(stringer))
}

func TestRedactedWithMultiByteValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "single rune", value: "é", want: "*"},
		{name: "multiple runes", value: "éa", want: "**"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, meta.Redacted(test.value).String())
		})
	}
}
