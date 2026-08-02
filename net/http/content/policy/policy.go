package policy

var undecodableKinds = map[string]struct{}{
	"gob":     {},
	"msgpack": {},
	"toml":    {},
}

// CanDecode reports whether kind is permitted to decode an untrusted HTTP request body.
//
// Unknown kinds are a separate concern: callers must first establish that a codec exists. This policy
// only rejects registered codecs whose decoders are neither ratio-bounded nor depth-bounded.
func CanDecode(kind string) bool {
	_, found := undecodableKinds[kind]
	return !found
}
