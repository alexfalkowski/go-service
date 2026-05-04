package meta

// NewPair creates one metadata key/value pair for batched storage updates.
//
// Use NewPair for custom metadata keys. Prefer the typed With* helpers in
// attrs.go for standard metadata keys so call sites stay consistent.
func NewPair(key string, value Value) Pair {
	return Pair{Key: key, Value: value}
}

// Pair stores one metadata key/value pair for batched storage updates.
//
// Pairs are passed to WithAttributes so multiple metadata attributes can be
// applied with one copy-on-write storage update.
type Pair struct {
	Key   string
	Value Value
}
