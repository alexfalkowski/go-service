package config

// Map for config.
type Map map[string]any

// Map at the key.
func (m Map) Map(key string) Map {
	return m[key].(Map)
}
