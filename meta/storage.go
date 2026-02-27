package meta

import (
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
// It is typically treated as immutable by callers (mutation happens through WithAttribute which
// returns a derived context containing an updated storage).
type Storage map[string]Value

// Add stores value under key and returns the Storage.
//
// Note: Storage is a map, so Add mutates the receiver.
func (s Storage) Add(key string, value Value) Storage {
	s[key] = value
	return s
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

func (s Storage) key(prefix, key string) string {
	if strings.IsEmpty(prefix) {
		return key
	}

	return prefix + key
}

func attributes(ctx context.Context) Storage {
	m := ctx.Value(meta)
	if m == nil {
		return make(Storage)
	}

	return m.(Storage)
}
