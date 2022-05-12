package config

import (
	"gopkg.in/yaml.v3"
)

// Map for config.
type Map map[string]any

// Map at the key.
func (m Map) Map(key string) Map {
	return m[key].(Map)
}

// UnmarshalFromBytes to map.
func UnmarshalFromBytes(bytes []byte, cfg Map) error {
	return yaml.Unmarshal(bytes, cfg)
}

// MarshalToBytes the map.
func MarshalToBytes(cfg Map) ([]byte, error) {
	return yaml.Marshal(cfg)
}
