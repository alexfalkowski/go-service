package content_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestNewStreamFromContentTypeRejectsEncoderOnlyKind(t *testing.T) {
	sm := stream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = nil
	sm.Register("json", codec)
	cont := content.NewContent(test.Encoder, sm, test.Pool)

	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.NDJSON)

	_, err := cont.NewStreamFromContentType(req)

	require.ErrorIs(t, err, content.ErrUnsupportedStreamMedia)
}

func TestNewStreamFromContentTypeResolvesNDJSON(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.NDJSON)

	m, err := test.Content.NewStreamFromContentType(req)

	require.NoError(t, err)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

func TestNewStreamFromMediaResolvesUsableDecoderForEveryStreamKind(t *testing.T) {
	// The client's response-stream path (see net/http/client/stream.go) relies on NewStreamFromMedia
	// staying unpoliced: it must keep resolving a usable decoder for every registered streaming kind,
	// including the ones NewStreamFromContentType rejects for untrusted request bodies.
	m, err := test.Content.NewStreamFromMedia(media.NDJSON)

	require.NoError(t, err)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

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
