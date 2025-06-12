package cli

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// FS for cli.
	FS = os.NewFS()

	// Name for cli.
	Name = env.NewName(FS)

	// Version for cli.
	Version = env.NewVersion()
)

func provide() (*os.FS, env.Name, env.Version) {
	return FS, Name, Version
}
