package cli

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	// FS is the filesystem used by CLI helpers for configuration lookup and reading sources.
	//
	// It defaults to an `*os.FS` rooted in the host filesystem (see `os.NewFS`). Tests may override
	// this variable to control filesystem reads performed during CLI setup.
	FS = os.NewFS()

	// Name is the CLI application name derived from the environment.
	//
	// The name is resolved via `env.NewName(FS)` and is used by `Application.Run` to populate the
	// command runner's app metadata (for example help text and descriptions).
	Name = env.NewName(FS)

	// Version is the CLI application version derived from the environment.
	//
	// The version is resolved via `env.NewVersion()` and is used by `Application.Run` to populate the
	// command runner's version information.
	Version = env.NewVersion()
)

func provide() (*os.FS, env.Name, env.Version) {
	return FS, Name, Version
}
