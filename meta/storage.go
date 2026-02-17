package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Converter transforms a key string before it is exported.
type Converter func(string) string

// Storage stores meta values keyed by attribute name.
type Storage map[string]Value

// Add a value with key.
func (s Storage) Add(key string, value Value) Storage {
	s[key] = value
	return s
}

// Get a value by key.
func (s Storage) Get(key string) Value {
	return s[key]
}

// Strings will create a map that is converts the key.
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
