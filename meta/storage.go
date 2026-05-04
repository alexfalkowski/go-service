package meta

import (
	"maps"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Converter transforms an attribute key before it is exported.
//
// It is used by export helpers to normalize key casing (for example snake_case or lowerCamelCase).
type Converter func(string) string

// Storage stores meta values keyed by attribute name.
//
// Storage is the internal backing map used by this package to hold context-scoped attributes.
// It is typically treated as immutable by callers. WithAttributes uses copy-on-write updates so
// derived contexts do not mutate parent storage.
type Storage map[string]Value

// AddPairs stores pairs in a copied Storage and returns the updated copy.
//
// The receiver is cloned once before all pairs are applied. This preserves context isolation while
// avoiding repeated map copies when several attributes are added together.
func (s Storage) AddPairs(pairs ...Pair) Storage {
	cloned := make(Storage, len(s)+len(pairs))
	maps.Copy(cloned, s)

	for _, pair := range pairs {
		cloned[pair.Key] = pair.Value
	}

	return cloned
}

// Get returns the value stored under key.
//
// If key is not present, Get returns the zero-value Value.
func (s Storage) Get(key string) Value {
	return s[key]
}

// Strings exports stored attributes as a string map.
//
// Each key is transformed using converter and then prefixed with prefix (if non-empty).
// Each value is rendered using Value.String().
//
// Export behavior:
//   - Attributes whose rendered value is an empty string are skipped.
//     (This includes Blank and Ignored values, and any Value whose rendered String() is empty.)
//   - Keys are included only if they have a non-empty rendered value.
func (s Storage) Strings(prefix string, converter Converter) Map {
	attributes := make(Map, len(s))
	for k, v := range s {
		if v := v.String(); !strings.IsEmpty(v) {
			attributes[s.key(prefix, converter(k))] = v
		}
	}
	return attributes
}

// Clone returns a shallow copy of the storage.
//
// Values are copied by assignment, which is sufficient because Value is immutable.
func (s Storage) Clone() Storage {
	cloned := make(Storage, len(s))
	maps.Copy(cloned, s)

	return cloned
}

func (s Storage) key(prefix, key string) string {
	if strings.IsEmpty(prefix) {
		return key
	}

	return prefix + key
}

func attributes(ctx context.Context) Storage {
	m := ctx.Value(meta)
	if m == nil {
		// Nil Storage avoids allocating on read-only paths; map reads, ranges, and the first AddPairs call are nil-safe.
		return nil
	}

	return m.(Storage)
}
