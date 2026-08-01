// Package policy defines HTTP content rules shared by unary and streaming request decoding.
package policy

const (
	gobKind         = "gob"
	messagePackKind = "msgpack"
)

// CanDecode reports whether kind is permitted to decode an untrusted HTTP request body.
//
// Unknown kinds are a separate concern: callers must first establish that a codec exists. This policy
// only rejects registered codecs whose decoders are neither ratio-bounded nor depth-bounded.
func CanDecode(kind string) bool {
	return kind != gobKind && kind != messagePackKind
}
