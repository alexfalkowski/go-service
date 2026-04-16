package time_test

import (
	"encoding/json"
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestMustParseDuration(t *testing.T) {
	require.Panics(t, func() { time.MustParseDuration("test") })
}

func TestDurationTextRoundTrip(t *testing.T) {
	duration := 5*time.Second + 250*time.Millisecond

	text, err := duration.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "5.25s", string(text))

	var decoded time.Duration
	require.NoError(t, decoded.UnmarshalText(text))
	require.Equal(t, duration, decoded)
}

func TestDurationJSONRoundTrip(t *testing.T) {
	duration := 3*time.Minute + 15*time.Second

	data, err := json.Marshal(duration)
	require.NoError(t, err)
	require.Equal(t, `"3m15s"`, string(data))

	var decoded time.Duration
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, duration, decoded)
}

func TestDurationUnmarshalTextInvalid(t *testing.T) {
	var duration time.Duration

	err := duration.UnmarshalText([]byte("not-a-duration"))
	require.Error(t, err)
}

func TestDurationUnmarshalJSONInvalid(t *testing.T) {
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
			var duration time.Duration

			err := json.Unmarshal([]byte(tt.input), &duration)
			require.Error(t, err)
		})
	}
}

func TestDurationZeroValueEncoding(t *testing.T) {
	var duration time.Duration

	text, err := duration.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "0s", string(text))

	data, err := json.Marshal(duration)
	require.NoError(t, err)
	require.Equal(t, `"0s"`, string(data))
}
