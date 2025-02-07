package meta

import (
	"context"
)

// Storage stores all the values for meta.
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
		if v.IsBlank() {
			continue
		}

		attributes[s.key(prefix, converter(k))] = v.String()
	}

	return attributes
}

func (s Storage) key(prefix, key string) string {
	if prefix == "" {
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
