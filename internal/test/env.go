package test

import (
	"github.com/alexfalkowski/go-service/env"
)

const (
	// Environment for test.
	Environment = env.Environment("dev")

	// Version for test.
	Version = env.Version("1.0.0")

	// Name for test.
	Name = env.Name("test")
)

// UserAgent for test.
var UserAgent = env.UserAgent(Name.String() + "/" + Version.String())
