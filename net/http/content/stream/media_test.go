package stream_test

import (
	"net/http/httptest"
	"testing"

	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestNewFromContentTypeRejectsEncoderOnlyKind(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Decoder = nil
	sm.Register("json", codec)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.NDJSON)

	_, err := contentstream.NewMediaFromContentType(req, sm)

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewFromAcceptRejectsDecoderOnlyKind(t *testing.T) {
	sm := encodingstream.NewMap()
	codec := sm.Get("json")
	codec.Encoder = nil
	sm.Register("json", codec)
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/hello", nil)
	req.Header.Set(content.AcceptKey, media.NDJSON)

	_, err := contentstream.NewMediaFromAccept(req, sm)

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewFromContentTypeResolvesNDJSON(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, media.NDJSON)

	m, err := contentstream.NewMediaFromContentType(req, test.StreamEncoder)

	require.NoError(t, err)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

func TestNewFromContentTypeRejectsWrongMajorNDJSON(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(content.TypeKey, "text/x-ndjson")

	_, err := contentstream.NewMediaFromContentType(req, test.StreamEncoder)

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewMediaResolvesUsableDecoderForEveryStreamKind(t *testing.T) {
	// The client's response-stream path (see net/http/client/stream.go) relies on NewMedia staying
	// unpoliced: it must keep resolving a usable decoder for every registered streaming kind,
	// including the ones NewFromContentType rejects for untrusted request bodies.
	m, err := contentstream.NewMedia(media.NDJSON, test.StreamEncoder)

	require.NoError(t, err)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

func TestNewMediaResolvesNDJSON(t *testing.T) {
	m, err := contentstream.NewMedia(media.NDJSON, encodingstream.NewMap())

	require.NoError(t, err)
	require.NotNil(t, m.NewEncoder)
	require.NotNil(t, m.NewDecoder)
	require.Equal(t, media.NDJSON, m.String())
}

func TestNewMediaRejectsUnmappedSubtype(t *testing.T) {
	_, err := contentstream.NewMedia(media.JSON, encodingstream.NewMap())

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewMediaRejectsWrongMajorNDJSON(t *testing.T) {
	_, err := contentstream.NewMedia("text/x-ndjson", encodingstream.NewMap())

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewMediaRejectsUnparsableType(t *testing.T) {
	_, err := contentstream.NewMedia("not-a-media-type", encodingstream.NewMap())

	require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
}

func TestNewFromAcceptResolvesWildcardOrExactMatchAnywhereInList(t *testing.T) {
	tests := []struct {
		name        string
		accept      string
		contentType string
	}{
		{name: "exact match", accept: media.NDJSON},
		{name: "bare wildcard", accept: "*/*"},
		{name: "bare wildcard with ndjson content-type", accept: "*/*", contentType: media.NDJSON},
		{name: "subtype wildcard", accept: "application/*"},
		{name: "exact match first in list", accept: media.NDJSON + ", */*"},
		{name: "wildcard first in list", accept: "*/*, " + media.NDJSON},
		{name: "browser-style list with trailing wildcard", accept: "text/html,application/xhtml+xml,*/*;q=0.8"},
		{name: "absent accept, ndjson content-type", contentType: media.NDJSON},
		{name: "absent accept, absent content-type"},
		{name: "exact match overrides a co-present zero quality bare wildcard", accept: "*/*;q=0, " + media.NDJSON},
		{name: "subtype wildcard overrides a co-present zero quality bare wildcard", accept: "*/*;q=0, application/*"},
		{name: "subtype wildcard overrides wrong-major zero quality type", accept: "text/x-ndjson;q=0, application/*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "GET", "/hello", nil)
			if tt.accept != "" {
				req.Header.Set(content.AcceptKey, tt.accept)
			}

			if tt.contentType != "" {
				req.Header.Set(content.TypeKey, tt.contentType)
			}

			m, err := contentstream.NewMediaFromAccept(req, test.StreamEncoder)

			require.NoError(t, err)
			require.NotNil(t, m.NewEncoder)
			require.Equal(t, media.NDJSON, m.String())
		})
	}
}

func TestNewFromAcceptRejectsConcreteOnlyUnsatisfiableAccept(t *testing.T) {
	tests := []struct {
		name   string
		accept string
	}{
		{name: "concrete unproducible type", accept: media.JSON},
		{name: "concrete type with wrong major", accept: "text/x-ndjson"},
		{name: "concrete list with no match and no wildcard", accept: media.JSON + ", " + media.YAML},
		{name: "unparsable", accept: "not-a-media-type"},
		{name: "subtype wildcard with mismatched major type", accept: "text/*"},
		{name: "bare wildcard with zero quality", accept: "*/*;q=0"},
		{name: "exact match with zero quality", accept: media.NDJSON + ";q=0"},
		{name: "zero quality wildcard alongside unproducible concrete type", accept: media.JSON + ", */*;q=0"},
		{name: "zero quality exact match overrides a co-present bare wildcard", accept: media.NDJSON + ";q=0, */*"},
		{name: "zero quality exact match overrides a co-present bare wildcard, reversed", accept: "*/*, " + media.NDJSON + ";q=0"},
		{name: "zero quality subtype wildcard overrides a co-present bare wildcard", accept: "application/*;q=0, */*"},
		{name: "zero quality subtype wildcard overrides a co-present bare wildcard, reversed", accept: "*/*, application/*;q=0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "GET", "/hello", nil)
			req.Header.Set(content.AcceptKey, tt.accept)

			_, err := contentstream.NewMediaFromAccept(req, test.StreamEncoder)

			require.ErrorIs(t, err, contentstream.ErrUnsupportedMedia)
		})
	}
}
