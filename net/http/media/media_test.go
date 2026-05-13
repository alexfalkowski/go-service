package media_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestParseInvalidMediaType(t *testing.T) {
	_, _, err := media.Parse("json")

	require.ErrorIs(t, err, media.ErrInvalidType)
	require.True(t, errors.Is(err, media.ErrInvalidType))
}

func TestWithUTF8(t *testing.T) {
	for _, tc := range []struct {
		name      string
		mediaType string
		expected  string
	}{
		{name: "text", mediaType: media.Text, expected: "text/plain; charset=utf-8"},
		{name: "html", mediaType: media.HTML, expected: "text/html; charset=utf-8"},
		{name: "existing charset", mediaType: "text/plain; charset=utf-8", expected: "text/plain; charset=utf-8"},
		{name: "parameter", mediaType: "text/plain; profile=test", expected: "text/plain; profile=test; charset=utf-8"},
		{name: "json", mediaType: media.JSON, expected: media.JSON},
		{name: "invalid", mediaType: "json", expected: "json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, media.WithUTF8(tc.mediaType))
		})
	}
}
