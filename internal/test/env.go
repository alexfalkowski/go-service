package test

import "github.com/alexfalkowski/go-service/env"

const (
	// ID for test.
	ID = env.ID("1234567890")

	// Name for test.
	Name = env.Name("test")

	// Version for test.
	Version = env.Version("1.0.0")

	// Environment for test.
	Environment = env.Environment("dev")
)

// UserAgent for test.
var UserAgent = env.UserAgent(Name.String() + "/" + Version.String())
