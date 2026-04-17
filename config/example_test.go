package config_test

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/os"
)

type exampleConfig struct {
	Name string `yaml:"name" validate:"required"`
}

func ExampleNewConfig() {
	fs := os.NewFS()

	dir, err := fs.MkdirTemp("", "go-service-config-example")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fs.RemoveAll(dir); err != nil {
			panic(err)
		}
	}()

	path := fs.Join(dir, "config.yml")
	if err := fs.WriteFile(path, []byte("name: payments"), 0o600); err != nil {
		panic(err)
	}

	flags := flag.NewFlagSet("serve")
	flags.AddInput("file:" + path)

	decoder := config.NewDecoder(config.DecoderParams{
		Flags:   flags,
		Encoder: exampleEncodingMap(),
		FS:      fs,
		Name:    env.Name("payments"),
	})

	cfg, err := config.NewConfig[exampleConfig](decoder, config.NewValidator())
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Name)
	// Output: payments
}

func exampleEncodingMap() *encoding.Map {
	return encoding.NewMap(encoding.MapParams{
		YAML: yaml.NewEncoder(),
	})
}
