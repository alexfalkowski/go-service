package mime

const (
	// ErrorMediaType is the media type used for plain-text error bodies.
	//
	// This is intended for responses where the body is a human-readable error message.
	// Note: "text/error" is not a standard IANA media type, but is used within go-service
	// for consistent internal error rendering.
	ErrorMediaType = "text/error; charset=utf-8"

	// HTMLMediaType is the media type for HTML documents encoded as UTF-8.
	//
	// This is typically used for HTML responses or debug pages.
	HTMLMediaType = "text/html; charset=utf-8"

	// JPEGMediaType is the media type for JPEG images.
	JPEGMediaType = "image/jpeg"

	// JSONMediaType is the media type for JSON documents.
	//
	// This is commonly used as the Content-Type for JSON request/response bodies.
	JSONMediaType = "application/json"

	// MarkdownMediaType is the media type for Markdown documents encoded as UTF-8.
	MarkdownMediaType = "text/markdown; charset=utf-8"

	// ProtobufMediaType is the media type for protobuf binary payloads.
	//
	// This is commonly used when transporting protobuf wire-format bodies over HTTP.
	ProtobufMediaType = "application/protobuf"

	// ProtobufJSONMediaType is the media type for protobuf JSON-encoded payloads.
	//
	// Note: this is a go-service specific media type string used to distinguish protobuf JSON
	// from generic JSON in content negotiation.
	ProtobufJSONMediaType = "application/pbjson"

	// ProtobufTextMediaType is the media type for protobuf text-format payloads.
	//
	// Note: this is a go-service specific media type string used to distinguish protobuf text format
	// in content negotiation.
	ProtobufTextMediaType = "application/pbtxt"

	// TextMediaType is the media type for plain text encoded as UTF-8.
	TextMediaType = "text/plain; charset=utf-8"

	// TOMLMediaType is the media type for TOML documents.
	TOMLMediaType = "application/toml"

	// YAMLMediaType is the media type for YAML documents.
	YAMLMediaType = "application/yaml"
)
