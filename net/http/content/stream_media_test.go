package content_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestNewStreamMediaResolvesNDJSON(t *testing.T) {
	m, err := content.NewStreamMedia(media.NDJSON, stream.NewMap())

	require.NoError(t, err)
	require.NotNil(t, m.NewEncoder)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

func TestNewStreamMediaRejectsUnmappedSubtype(t *testing.T) {
	_, err := content.NewStreamMedia(media.JSON, stream.NewMap())

	require.ErrorIs(t, err, content.ErrUnsupportedStreamMedia)
}

func TestNewStreamMediaRejectsUnparsableType(t *testing.T) {
	_, err := content.NewStreamMedia("not-a-media-type", stream.NewMap())

	require.ErrorIs(t, err, content.ErrUnsupportedStreamMedia)
}
