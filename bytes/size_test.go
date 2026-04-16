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

func TestSizeTextRoundTrip(t *testing.T) {
	size := bytes.MustParseSize("4MB")

	text, err := size.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "4MB", string(text))

	var decoded bytes.Size
	require.NoError(t, decoded.UnmarshalText(text))
	require.Equal(t, size, decoded)
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
