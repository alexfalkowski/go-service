package ristretto

// Config for ristretto.
type Config struct {
	Name        string `yaml:"name"`
	NumCounters int64  `yaml:"max_counters"`
	MaxCost     int64  `yaml:"max_cost"`
	BufferItems int64  `yaml:"buffer_items"`
}
