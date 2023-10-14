package ristretto

// Config for ristretto.
type Config struct {
	NumCounters int64 `yaml:"num_counters" json:"num_counters" toml:"num_counters"`
	MaxCost     int64 `yaml:"max_cost" json:"max_cost" toml:"max_cost"`
	BufferItems int64 `yaml:"buffer_items" json:"buffer_items" toml:"buffer_items"`
}
