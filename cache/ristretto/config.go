package ristretto

// Config for ristretto.
type Config struct {
	NumCounters int64 `yaml:"num_counters,omitempty" json:"num_counters,omitempty" toml:"num_counters,omitempty"`
	MaxCost     int64 `yaml:"max_cost,omitempty" json:"max_cost,omitempty" toml:"max_cost,omitempty"`
	BufferItems int64 `yaml:"buffer_items,omitempty" json:"buffer_items,omitempty" toml:"buffer_items,omitempty"`
}
