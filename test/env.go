package test

import (
	"github.com/alexfalkowski/go-service/env"
)

var (
	// Environment for test.
	Environment = env.Environment("dev")

	// Version for test.
	Version = env.Version("1.0.0")

	// Name for test.
	Name = env.Name("test")

	// UserAgent for test.
	UserAgent = env.UserAgent(string(Name) + "/" + string(Version))
)
