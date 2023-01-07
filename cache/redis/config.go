package redis

// Config for redis.
type Config struct {
	Addresses map[string]string `yaml:"addresses" json:"addresses" toml:"addresses"`
	Username  string            `yaml:"username" json:"username" toml:"username"`
	Password  string            `yaml:"password" json:"password" toml:"password"`
	DB        int               `yaml:"db" json:"db" toml:"db"`
}
