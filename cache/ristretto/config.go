package ristretto

// Config for ristretto.
type Config struct {
	NumCounters int64 `yaml:"max_counters"`
	MaxCost     int64 `yaml:"max_cost"`
	BufferItems int64 `yaml:"buffer_items"`
}
