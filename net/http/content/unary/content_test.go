package unary_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestNewFromMediaSelectsUsableCodec(t *testing.T) {
	for _, tc := range mediaTests() {
		t.Run(tc.name, func(t *testing.T) {
			media := test.UnaryContent.NewFromMedia(tc.mediaType)

			require.Equal(t, tc.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
		})
	}
}

func TestContentSelectionUsesRequestMediaType(t *testing.T) {
	for _, tc := range mediaTests() {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.ContentTypeKey, tc.mediaType)

			media := test.UnaryContent.NewFromRequest(req)

			require.Equal(t, tc.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
		})
	}
}

func TestNewFromRequestPrefersContentType(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.TOML)
	req.Header.Set(http.AcceptKey, media.YAML)

	media := test.UnaryContent.NewFromRequest(req)

	require.Equal(t, "toml", media.Subtype())
	require.Same(t, test.Encoder.Get("toml"), media.Encoder)
}

func TestNewFromAcceptPrefersAccept(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.TOML)
	req.Header.Set(http.AcceptKey, media.YAML)

	media := test.UnaryContent.NewFromAccept(req)

	require.Equal(t, "yaml", media.Subtype())
	require.Same(t, test.Encoder.Get("yaml"), media.Encoder)
}

func TestNewFromAcceptUsesFirstCompleteMediaRange(t *testing.T) {
	tests := []struct {
		name    string
		accept  string
		subtype string
		kind    string
	}{
		{name: "quoted comma", accept: `application/yaml; profile="a,b", application/toml`, subtype: "yaml", kind: "yaml"},
		{name: "escaped quoted comma", accept: `application/yaml; profile="a\",b", application/toml`, subtype: "yaml", kind: "yaml"},
		{name: "list", accept: "application/yaml, application/toml", subtype: "yaml", kind: "yaml"},
		{name: "malformed quoted parameter", accept: `application/yaml; profile="a,b`, subtype: "json", kind: "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.AcceptKey, tt.accept)

			media := test.UnaryContent.NewFromAccept(req)

			require.Equal(t, tt.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tt.kind), media.Encoder)
		})
	}
}

func TestNewFromAcceptFallsBackToContentType(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.TOML)

	media := test.UnaryContent.NewFromAccept(req)

	require.Equal(t, "toml", media.Subtype())
	require.Same(t, test.Encoder.Get("toml"), media.Encoder)
}

func TestNewFromRequestFallsBackToAccept(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.AcceptKey, `application/yaml; profile="a,b", application/toml`)

	media := test.UnaryContent.NewFromRequest(req)

	require.Equal(t, "yaml", media.Subtype())
	require.Same(t, test.Encoder.Get("yaml"), media.Encoder)
}

func TestNewFromRequestNormalizesMediaTypeCase(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, "Application/YAML")

	media := test.UnaryContent.NewFromRequest(req)

	require.Equal(t, "yaml", media.Subtype())
	require.Same(t, test.Encoder.Get("yaml"), media.Encoder)
}

func TestNewFromRequestNormalizesAcceptMediaTypeCase(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.AcceptKey, "Application/YAML, Application/TOML")

	media := test.UnaryContent.NewFromRequest(req)

	require.Equal(t, "yaml", media.Subtype())
	require.Same(t, test.Encoder.Get("yaml"), media.Encoder)
}

func TestNewFromRequestFallsBackFromInternalErrorMedia(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.AcceptKey, media.Error)

	media := test.UnaryContent.NewFromRequest(req)

	require.Equal(t, "plain", media.Subtype())
	require.Same(t, test.Encoder.Get("bytes"), media.Encoder)
}

func TestContentSelectionUsesContentType(t *testing.T) {
	for _, tc := range mediaTests() {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.ContentTypeKey, tc.mediaType)

			media := test.UnaryContent.NewFromContentType(req)

			require.Equal(t, tc.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
		})
	}
}

func TestNewFromContentTypeFallsBackFromInternalErrorMedia(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.Error)

	media := test.UnaryContent.NewFromContentType(req)

	require.Equal(t, "plain", media.Subtype())
	require.Same(t, test.Encoder.Get("bytes"), media.Encoder)
}

func TestNewFromRequestBodyRejectsUnsafeMedia(t *testing.T) {
	for _, mediaType := range []string{"application/gob", media.MessagePack, media.MessagePack + "; profile=test", media.TOML, media.TOML + "; profile=test"} {
		t.Run(mediaType, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.ContentTypeKey, mediaType)

			media, err := test.UnaryContent.NewFromRequestBody(req)

			require.ErrorIs(t, err, unary.ErrUnsupportedRequestMedia)
			require.EqualError(t, err, "unary: unsupported request media")
			require.False(t, media.CanDecodeRequest())
		})
	}
}

func TestNewFromRequestBodyRejectsUnknownContentType(t *testing.T) {
	for _, mediaType := range []string{"application/cbor", "application/yml", "/"} {
		t.Run(mediaType, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.ContentTypeKey, mediaType)

			media, err := test.UnaryContent.NewFromRequestBody(req)

			require.ErrorIs(t, err, unary.ErrUnsupportedRequestMedia)
			require.False(t, media.CanDecodeRequest())
		})
	}
}

func TestNewFromRequestBodyRejectsNilRegisteredCodec(t *testing.T) {
	tests := []struct {
		name        string
		register    string
		contentType string
	}{
		{name: "absent content type, json nil-registered", register: "json", contentType: ""},
		{name: "text/error, bytes nil-registered", register: "bytes", contentType: media.Error},
		{name: "msgpack, msgpack nil-registered", register: "msgpack", contentType: media.MessagePack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := encoding.NewMap()
			enc.Register(tt.register, nil)
			cont := unary.NewContent(enc, test.Pool)

			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			if tt.contentType != "" {
				req.Header.Set(http.ContentTypeKey, tt.contentType)
			}

			media, err := cont.NewFromRequestBody(req)

			require.ErrorIs(t, err, unary.ErrUnsupportedRequestMedia)
			require.False(t, media.CanDecodeRequest())
		})
	}
}

func TestNewFromRequestBodyDefaultsToJSONWhenContentTypeAbsent(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)

	media, err := test.UnaryContent.NewFromRequestBody(req)

	require.NoError(t, err)
	require.Equal(t, "json", media.Subtype())
	require.Same(t, test.Encoder.Get("json"), media.Encoder)
	require.True(t, media.CanDecodeRequest())
}

func TestNewFromRequestBodyTreatsParameterizedInternalErrorContentTypeAsText(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
	req.Header.Set(http.ContentTypeKey, media.Error+"; charset=utf-8")

	m, err := test.UnaryContent.NewFromRequestBody(req)

	require.NoError(t, err)
	require.Equal(t, "plain", m.Subtype())
	require.Same(t, test.Encoder.Get("bytes"), m.Encoder)
	require.True(t, m.CanDecodeRequest())
}

func TestNewFromRequestBodyDecodesFallthroughReachableMediaTypes(t *testing.T) {
	// Guards against the identity fallthrough (knownMedia -> media.Parse -> enc.Get(unaryKind(subtype)))
	// being reverted into an allowlist: these media types have no dedicated knownMedia case and must keep
	// resolving through the general parser path.
	for _, tc := range []struct {
		mediaType string
		subtype   string
		kind      string
	}{
		{mediaType: "application/pb", subtype: "pb", kind: "protobuf"},
		{mediaType: "application/protobin", subtype: "protobin", kind: "protobuf"},
		{mediaType: "application/octet-stream", subtype: "octet-stream", kind: "bytes"},
		{mediaType: media.Text, subtype: "plain", kind: "bytes"},
	} {
		t.Run(tc.mediaType, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			req.Header.Set(http.ContentTypeKey, tc.mediaType)

			media, err := test.UnaryContent.NewFromRequestBody(req)

			require.NoError(t, err)
			require.Equal(t, tc.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tc.kind), media.Encoder)
			require.True(t, media.CanDecodeRequest())
		})
	}
}

func TestNewFromAcceptResolvesJSONForWildcardOrUnknown(t *testing.T) {
	// Outbound (Accept) negotiation keeps its JSON fallback for an unproducible or absent preference;
	// strict Accept parsing (q-values, wildcards) is out of scope for this package.
	for _, tc := range []struct {
		name   string
		accept string
	}{
		{name: "wildcard", accept: "*/*"},
		{name: "unproducible", accept: media.HTML},
		{name: "absent", accept: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/hello", nil)
			if tc.accept != "" {
				req.Header.Set(http.AcceptKey, tc.accept)
			}

			m := test.UnaryContent.NewFromAccept(req)

			require.Equal(t, "json", m.Subtype())
			require.Same(t, test.Encoder.Get("json"), m.Encoder)
		})
	}
}

func TestEveryEncoderKindIsClassified(t *testing.T) {
	// Every kind in the default registry must be explicitly classified, so a new codec cannot become
	// request-decodable without a decision. See policy.CanDecode.
	classified := map[string]bool{
		"json": true, "hjson": true, "yaml": true, "toml": false, "bytes": true,
		"protobuf": true, "prototext": true, "protojson": true,
		"msgpack": false, "gob": false,
	}

	for _, kind := range encoding.NewMap().Keys() {
		expected, ok := classified[kind]
		require.True(t, ok, "kind %q needs an explicit request-decode classification", kind)
		content := unary.NewContent(encoding.NewMap(), test.Pool)
		require.Equal(t, expected, content.NewFromMedia("application/"+kind).CanDecodeRequest(), "kind %q", kind)
	}
}

func TestContentSelectionParsesParameterizedMediaType(t *testing.T) {
	tests := []struct {
		name      string
		mediaType string
		expected  string
		subtype   string
		kind      string
	}{
		{name: "json", mediaType: "application/json; profile=test", expected: media.JSON, subtype: "json", kind: "json"},
		{name: "msgpack", mediaType: "application/vnd.msgpack; profile=test", expected: media.MessagePack, subtype: "msgpack", kind: "msgpack"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := test.UnaryContent.NewFromMedia(tt.mediaType)

			require.Equal(t, tt.expected, media.String())
			require.Equal(t, tt.subtype, media.Subtype())
			require.Same(t, test.Encoder.Get(tt.kind), media.Encoder)
		})
	}
}

func TestNewFromMediaPreservesInternalErrorMedia(t *testing.T) {
	media := test.UnaryContent.NewFromMedia(media.Error)

	require.Equal(t, "error", media.Subtype())
	require.Nil(t, media.Encoder)
}

func TestNewFromMediaFallsBackToJSONWhenKnownEncoderIsMissing(t *testing.T) {
	enc := &encoding.Map{}
	media := unary.NewContent(enc, test.Pool).NewFromMedia(media.HumanJSON)

	require.Equal(t, "json", media.Subtype())
	require.Equal(t, "application/json", media.String())
	require.Nil(t, media.Encoder)
}

type mediaTest struct {
	name      string
	mediaType string
	subtype   string
	kind      string
}

func mediaTests() []mediaTest {
	return []mediaTest{
		{name: "json", mediaType: media.JSON, subtype: "json", kind: "json"},
		{name: "hjson", mediaType: media.HumanJSON, subtype: "hjson", kind: "hjson"},
		{name: "yaml", mediaType: media.YAML, subtype: "yaml", kind: "yaml"},
		{name: "toml", mediaType: media.TOML, subtype: "toml", kind: "toml"},
		{name: "msgpack", mediaType: media.MessagePack, subtype: "msgpack", kind: "msgpack"},
		{name: "protobuf", mediaType: media.Protobuf, subtype: "protobuf", kind: "protobuf"},
		{name: "proto", mediaType: "application/proto", subtype: "proto", kind: "protobuf"},
		{name: "pb", mediaType: "application/pb", subtype: "pb", kind: "protobuf"},
		{name: "protobin", mediaType: "application/protobin", subtype: "protobin", kind: "protobuf"},
		{name: "pbbin", mediaType: "application/pbbin", subtype: "pbbin", kind: "protobuf"},
		{name: "protobuf json", mediaType: media.ProtobufJSON, subtype: "pbjson", kind: "protojson"},
		{name: "protojson", mediaType: "application/protojson", subtype: "protojson", kind: "protojson"},
		{name: "protobuf text", mediaType: media.ProtobufText, subtype: "pbtxt", kind: "prototext"},
		{name: "prototext", mediaType: "application/prototext", subtype: "prototext", kind: "prototext"},
		{name: "prototxt", mediaType: "application/prototxt", subtype: "prototxt", kind: "prototext"},
		{name: "gob", mediaType: "application/gob", subtype: "gob", kind: "gob"},
		{name: "plain", mediaType: media.Text, subtype: "plain", kind: "bytes"},
		{name: "octet-stream", mediaType: "application/octet-stream", subtype: "octet-stream", kind: "bytes"},
		{name: "invalid", mediaType: "test", subtype: "json", kind: "json"},
		{name: "unknown", mediaType: "application/test", subtype: "json", kind: "json"},
	}
}
