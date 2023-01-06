package redis

// Config for redis.
type Config struct {
	Addresses map[string]string `yaml:"addresses" json:"addresses"`
	Username  string            `yaml:"username" json:"username"`
	Password  string            `yaml:"password" json:"password"`
	DB        int               `yaml:"db" json:"db"`
}
