package meta

import (
	"context"
)

type (
	// Converter takes a string and creates a new string.
	Converter func(string) string

	// Storage stores all the values for meta.
	Storage map[string]*Value
)

// Add a value with key.
func (s Storage) Add(key string, value *Value) Storage {
	s[key] = value

	return s
}

// Get a value by key.
func (s Storage) Get(key string) *Value {
	return s[key]
}

// Strings will create a map that is converts the key.
func (s Storage) Strings(prefix string, converter Converter) Map {
	attributes := make(Map, len(s))

	for k, v := range s {
		if v := v.String(); v != "" {
			attributes[s.key(prefix, converter(k))] = v
		}
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
