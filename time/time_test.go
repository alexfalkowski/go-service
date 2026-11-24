package time_test

import (
	"testing"
	"testing/synctest"

	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/stretchr/testify/require"
)

func TestParsePreservesTimestampPrecision(t *testing.T) {
	const layout = "02 Jan 2006 15:04:05.000000000 -0700"

	tm, err := time.Parse(layout, "24 Nov 2026 13:30:00.123456789 +0100")
	require.NoError(t, err)
	require.Equal(t, "2026-11-24T13:30:00.123456789+01:00", tm.Format(time.RFC3339Nano))
}

func TestParseReturnsParseErrorForInvalidInput(t *testing.T) {
	_, err := time.Parse(time.RFC3339, "invalid")

	var parseError *time.ParseError
	require.ErrorAs(t, err, &parseError)
}

func TestParseInLocationUsesLocationForUnzonedInput(t *testing.T) {
	loc := time.FixedZone("UTC+01", 60*60)

	tm, err := time.ParseInLocation(time.DateTime, "2026-11-24 13:30:00", loc)
	require.NoError(t, err)
	require.Equal(t, "2026-11-24T12:30:00Z", tm.UTC().Format(time.RFC3339))
}

func TestDateNormalizesFieldsInSpecifiedLocation(t *testing.T) {
	loc := time.FixedZone("UTC+05:30", 5*60*60+30*60)

	tm := time.Date(2026, time.January, 32, 6, 7, 8, 123456789, loc)

	require.Equal(t, "2026-02-01T06:07:08.123456789+05:30", tm.Format(time.RFC3339Nano))
	name, _ := tm.Zone()
	require.Equal(t, "UTC+05:30", name)
}

func TestLoadLocationResolvesUTCAndRejectsInvalidName(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	require.NoError(t, err)
	require.Equal(t, "UTC", loc.String())

	_, err = time.LoadLocation("../UTC")
	require.Error(t, err)
}

func TestLoadLocationFromTZDataUsesSuppliedZone(t *testing.T) {
	const zeroWord = "\x00\x00\x00\x00"
	const tzif = "TZif" +
		zeroWord + zeroWord + zeroWord + zeroWord + // Version 1 and reserved bytes.
		zeroWord + zeroWord + zeroWord + zeroWord + // No indicators, leap seconds, or transitions.
		"\x00\x00\x00\x01\x00\x00\x00\x04" + // One zone and four abbreviation bytes.
		"\x00\x00\x2a\x30\x00\x00" + // UTC+03:00, standard time, abbreviation at index zero.
		"TST\x00"

	loc, err := time.LoadLocationFromTZData("Test/Fixed", []byte(tzif))
	require.NoError(t, err)
	require.Equal(t, "Test/Fixed", loc.String())

	name, offset := time.Unix(0).In(loc).Zone()
	require.Equal(t, "TST", name)
	require.Equal(t, 3*60*60, offset)

	_, err = time.LoadLocationFromTZData("Test/Fixed", []byte("invalid"))
	require.Error(t, err)
}

func TestUnixSubsecondConstructorsPreservePrecisionBeforeEpoch(t *testing.T) {
	tests := []struct {
		construct func(int64) time.Time
		name      string
		want      string
	}{
		{construct: time.UnixMilli, name: "milliseconds", want: "1969-12-31T23:59:58.766Z"},
		{construct: time.UnixMicro, name: "microseconds", want: "1969-12-31T23:59:59.998766Z"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tm := test.construct(-1234)

			require.Equal(t, test.want, tm.UTC().Format(time.RFC3339Nano))
		})
	}
}

func TestAfterFuncCallsFunctionAtDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		called := make(chan time.Time, 1)

		timer := time.AfterFunc(time.Second, func() {
			called <- time.Now()
		})
		require.NotNil(t, timer)
		t.Cleanup(func() {
			timer.Stop()
		})

		require.Equal(t, time.Second.Duration(), (<-called).Sub(start))
	})
}

func TestTickDeliversRepeatedIntervals(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		ticks := time.Tick(time.Second)
		require.NotNil(t, ticks)

		first := <-ticks
		second := <-ticks

		require.Equal(t, time.Second.Duration(), first.Sub(start))
		require.Equal(t, time.Second.Duration(), second.Sub(first))
	})
}

func TestTickReturnsNilForNonPositiveIntervals(t *testing.T) {
	require.Nil(t, time.Tick(0))
	require.Nil(t, time.Tick(-time.Second))
}

func TestNewTimerDeliversConfiguredDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		timer := time.NewTimer(time.Nanosecond)
		require.NotNil(t, timer)
		t.Cleanup(func() {
			timer.Stop()
		})

		tm := <-timer.C
		require.False(t, tm.IsZero())
	})
}

func TestUntilReturnsDurationUntilFutureTime(t *testing.T) {
	future := time.Now().Add(time.Hour.Duration())

	duration := time.Until(future)

	require.Positive(t, duration)
	require.LessOrEqual(t, duration, time.Hour)
}
