package redis

// Config for redis.
type Config struct {
	Addresses map[string]string `yaml:"addresses,omitempty" json:"addresses,omitempty" toml:"addresses,omitempty"`
	Username  string            `yaml:"username,omitempty" json:"username,omitempty" toml:"username,omitempty"`
	Password  string            `yaml:"password,omitempty" json:"password,omitempty" toml:"password,omitempty"`
	DB        int               `yaml:"db,omitempty" json:"db,omitempty" toml:"db,omitempty"`
}
