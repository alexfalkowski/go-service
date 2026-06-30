package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

// FuzzDurationTextRoundTrip explores duration parsing and text marshaling used by config values.
func FuzzDurationTextRoundTrip(f *testing.F) {
	f.Add("0s")
	f.Add("250ms")
	f.Add("5.25s")
	f.Add("3m15s")
	f.Add("not-a-duration")

	f.Fuzz(func(t *testing.T, text string) {
		if len(text) > 128 {
			t.Skip()
		}

		var duration time.Duration
		if err := duration.UnmarshalText([]byte(text)); err != nil {
			return
		}

		encoded, err := duration.MarshalText()
		require.NoError(t, err)

		var decoded time.Duration
		require.NoError(t, decoded.UnmarshalText(encoded))
		require.Equal(t, duration, decoded)
	})
}

// FuzzDurationJSONRoundTrip explores JSON duration parsing while preserving the text round-trip invariant.
func FuzzDurationJSONRoundTrip(f *testing.F) {
	f.Add([]byte(`"0s"`))
	f.Add([]byte(`"250ms"`))
	f.Add([]byte(`"5.25s"`))
	f.Add([]byte(`"3m15s"`))
	f.Add([]byte(`"not-a-duration"`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1024 {
			t.Skip()
		}

		var duration time.Duration
		if err := json.Unmarshal(data, &duration); err != nil {
			return
		}

		encoded, err := json.Marshal(duration)
		require.NoError(t, err)

		var decoded time.Duration
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		require.Equal(t, duration, decoded)
	})
}
