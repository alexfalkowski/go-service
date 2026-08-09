package meta_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestValueFormatsVisibleIgnoredAndRedactedValues(t *testing.T) {
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

func TestErrorReturnsVisibleErrorMessage(t *testing.T) {
	value := meta.Error(&test.MessageError{Message: "boom"})

	require.Equal(t, "boom", value.Value())
	require.Equal(t, "boom", value.String())
	require.False(t, value.IsEmpty())
}

func TestErrorReturnsBlankValueForNilErrors(t *testing.T) {
	var typedNil error = (*test.MessageError)(nil)

	for _, test := range []struct {
		name string
		err  error
	}{
		{name: "nil"},
		{name: "typed nil", err: typedNil},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, meta.Blank(), meta.Error(test.err))
		})
	}
}

func TestValueConversionsUseConfiguredStringers(t *testing.T) {
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

func TestValueConversionsHandleNilStringers(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	tests := []struct {
		name    string
		convert func(fmt.Stringer) meta.Value
		value   fmt.Stringer
	}{
		{name: "to string with nil", convert: meta.ToString},
		{name: "to string with typed nil", convert: meta.ToString, value: stringer},
		{name: "to redacted with typed nil", convert: meta.ToRedacted, value: stringer},
		{name: "to ignored with typed nil", convert: meta.ToIgnored, value: stringer},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, meta.Blank(), test.convert(test.value))
		})
	}
}

func TestRedactedValuePreservesMultiByteCharacters(t *testing.T) {
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
