package cli

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// FS is the filesystem used for CLI configuration lookup.
	FS = os.NewFS()

	// Name is the CLI application name derived from the environment.
	Name = env.NewName(FS)

	// Version is the CLI application version derived from the environment.
	Version = env.NewVersion()
)

func provide() (*os.FS, env.Name, env.Version) {
	return FS, Name, Version
}
