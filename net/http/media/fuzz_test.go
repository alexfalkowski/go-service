package media_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

// FuzzParse explores media type parsing and normalization for request content negotiation.
func FuzzParse(f *testing.F) {
	for _, value := range []string{
		"json",
		"text/plain; charset",
		"Application/YAML; profile=test",
		media.MessagePack,
		"TEXT/PLAIN; Charset=utf-8",
		media.Text,
		media.HTML,
		"text/plain; profile=test",
		media.JSON,
	} {
		f.Add(value)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if len(value) > 512 {
			t.Skip()
		}

		mediaType, err := media.Parse(value)
		if err != nil {
			require.ErrorIs(t, err, media.ErrInvalidType)
			return
		}

		require.False(t, mediaType.IsZero())
		require.NotEmpty(t, mediaType.String())
		require.NotEmpty(t, mediaType.Subtype())
		require.NotEmpty(t, mediaType.WithUTF8())
	})
}
