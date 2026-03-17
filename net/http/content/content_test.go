package content_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/stretchr/testify/require"
)

func TestNewFromMedia(t *testing.T) {
	for _, tc := range mediaTests() {
		t.Run(tc.name, func(t *testing.T) {
			media := test.Content.NewFromMedia(tc.mediaType)

			require.Equal(t, tc.subtype, media.Subtype)
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
		})
	}
}

func TestNewFromRequest(t *testing.T) {
	for _, tc := range mediaTests() {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(content.TypeKey, tc.mediaType)

			media := test.Content.NewFromRequest(req)

			require.Equal(t, tc.subtype, media.Subtype)
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
		})
	}
}

type mediaTest struct {
	name      string
	mediaType string
	subtype   string
	kind      string
}

func mediaTests() []mediaTest {
	return []mediaTest{
		{name: "json", mediaType: mime.JSONMediaType, subtype: "json", kind: "json"},
		{name: "hjson", mediaType: mime.HJSONMediaType, subtype: "hjson", kind: "hjson"},
		{name: "yaml", mediaType: mime.YAMLMediaType, subtype: "yaml", kind: "yaml"},
		{name: "yml", mediaType: "application/yml", subtype: "yml", kind: "yml"},
		{name: "toml", mediaType: mime.TOMLMediaType, subtype: "toml", kind: "toml"},
		{name: "protobuf", mediaType: mime.ProtobufMediaType, subtype: "protobuf", kind: "protobuf"},
		{name: "proto", mediaType: "application/proto", subtype: "proto", kind: "proto"},
		{name: "pb", mediaType: "application/pb", subtype: "pb", kind: "pb"},
		{name: "protobin", mediaType: "application/protobin", subtype: "protobin", kind: "protobin"},
		{name: "pbbin", mediaType: "application/pbbin", subtype: "pbbin", kind: "pbbin"},
		{name: "protobuf json", mediaType: mime.ProtobufJSONMediaType, subtype: "pbjson", kind: "pbjson"},
		{name: "protojson", mediaType: "application/protojson", subtype: "protojson", kind: "protojson"},
		{name: "protobuf text", mediaType: mime.ProtobufTextMediaType, subtype: "pbtxt", kind: "pbtxt"},
		{name: "prototext", mediaType: "application/prototext", subtype: "prototext", kind: "prototext"},
		{name: "prototxt", mediaType: "application/prototxt", subtype: "prototxt", kind: "prototxt"},
		{name: "gob", mediaType: "application/gob", subtype: "gob", kind: "gob"},
		{name: "plain", mediaType: mime.TextMediaType, subtype: "plain", kind: "plain"},
		{name: "octet-stream", mediaType: "application/octet-stream", subtype: "octet-stream", kind: "octet-stream"},
		{name: "markdown", mediaType: mime.MarkdownMediaType, subtype: "markdown", kind: "markdown"},
		{name: "invalid", mediaType: "test", subtype: "json", kind: "json"},
		{name: "unknown", mediaType: "application/test", subtype: "json", kind: "json"},
	}
}
