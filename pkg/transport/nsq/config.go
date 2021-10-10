package nsq

type Config struct {
	LookupHost string `yaml:"lookup_host"`
	Host       string `yaml:"host"`
}
