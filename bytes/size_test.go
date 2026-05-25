package bytes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/stretchr/testify/require"
)

func TestMustParseSize(t *testing.T) {
	require.Panics(t, func() { bytes.MustParseSize("test") })
}

func TestDefaultSize(t *testing.T) {
	require.Equal(t, 4*bytes.MB, bytes.DefaultSize)
}

func TestSizeTextRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		text string
		size bytes.Size
	}{
		{name: "decimal boundary", size: bytes.MustParseSize("4MB"), text: "4000000B"},
		{name: "above byte boundary", size: bytes.Size(1001), text: "1001B"},
		{name: "binary megabyte", size: bytes.Size(1048576), text: "1048576B"},
		{name: "below decimal megabyte", size: bytes.Size(999999), text: "999999B"},
		{name: "above decimal megabyte", size: bytes.Size(1000001), text: "1000001B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := tt.size.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tt.text, string(text))

			var decoded bytes.Size
			require.NoError(t, decoded.UnmarshalText(text))
			require.Equal(t, tt.size, decoded)
		})
	}
}

func TestSizeJSONRoundTrip(t *testing.T) {
	size := bytes.MustParseSize("64B")

	data, err := json.Marshal(size)
	require.NoError(t, err)
	require.Equal(t, `"64B"`, string(data))

	var decoded bytes.Size
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, size, decoded)
}

func TestSizeJSONRoundTripPreservesByteCount(t *testing.T) {
	size := bytes.Size(1048576)

	data, err := json.Marshal(size)
	require.NoError(t, err)
	require.Equal(t, `"1048576B"`, string(data))

	var decoded bytes.Size
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, size, decoded)
}

func TestSizeUnmarshalTextInvalid(t *testing.T) {
	var size bytes.Size

	err := size.UnmarshalText([]byte("not-a-size"))
	require.Error(t, err)
}

func TestSizeUnmarshalJSONInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "null", input: "null"},
		{name: "number", input: "5"},
		{name: "object", input: "{}"},
		{name: "invalid string value", input: `"bad"`},
		{name: "malformed string", input: `"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var size bytes.Size

			err := json.Unmarshal([]byte(tt.input), &size)
			require.Error(t, err)
		})
	}
}

func TestSizeZeroValueEncoding(t *testing.T) {
	var size bytes.Size

	text, err := size.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "0B", string(text))

	data, err := json.Marshal(size)
	require.NoError(t, err)
	require.Equal(t, `"0B"`, string(data))
}
