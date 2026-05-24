package media_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestParseInvalidMediaType(t *testing.T) {
	_, err := media.Parse("json")
	require.ErrorIs(t, err, media.ErrInvalidType)
}

func TestParseInvalidMediaTypeParameter(t *testing.T) {
	_, err := media.Parse("text/plain; charset")
	require.ErrorIs(t, err, media.ErrInvalidType)
}

func TestParseNormalizesMediaTypeCase(t *testing.T) {
	mediaType, err := media.Parse("Application/YAML; profile=test")

	require.NoError(t, err)
	require.False(t, mediaType.IsZero())
	require.Equal(t, media.YAML, mediaType.String())
	require.Equal(t, "yaml", mediaType.Subtype())
}

func TestParseStripsVendorSubtypePrefix(t *testing.T) {
	mediaType, err := media.Parse(media.MessagePack)

	require.NoError(t, err)
	require.Equal(t, media.MessagePack, mediaType.String())
	require.Equal(t, "msgpack", mediaType.Subtype())
}

func TestTypeWithUTF8(t *testing.T) {
	mediaType, err := media.Parse("TEXT/PLAIN; Charset=utf-8")

	require.NoError(t, err)
	require.Equal(t, "TEXT/PLAIN; Charset=utf-8", mediaType.WithUTF8())
}

func TestTypeByExtension(t *testing.T) {
	require.Equal(t, "image/svg+xml", media.TypeByExtension(".svg"))
}

func TestTypeWithUTF8Values(t *testing.T) {
	for _, tc := range []struct {
		name      string
		mediaType string
		expected  string
	}{
		{name: "text", mediaType: media.Text, expected: "text/plain; charset=utf-8"},
		{name: "html", mediaType: media.HTML, expected: "text/html; charset=utf-8"},
		{name: "existing charset", mediaType: "text/plain; charset=utf-8", expected: "text/plain; charset=utf-8"},
		{name: "existing charset with case variant", mediaType: "text/plain; Charset=utf-8", expected: "text/plain; Charset=utf-8"},
		{name: "parameter", mediaType: "text/plain; profile=test", expected: "text/plain; profile=test; charset=utf-8"},
		{name: "json", mediaType: media.JSON, expected: media.JSON},
		{name: "msgpack", mediaType: media.MessagePack, expected: media.MessagePack},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mediaType, err := media.Parse(tc.mediaType)

			require.NoError(t, err)
			require.Equal(t, tc.expected, mediaType.WithUTF8())
		})
	}
}
