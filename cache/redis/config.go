package redis

// Config for redis.
type Config struct {
	Addresses map[string]string `yaml:"addresses"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	DB        int               `yaml:"db"`
}
