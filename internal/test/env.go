package test

import (
	"github.com/alexfalkowski/go-service/v2/env"
)

const (
	// ID for test.
	ID = env.ID("1234567890")

	// Name for test.
	Name = env.Name("test")

	// UserID for test.
	UserID = env.UserID(Name)

	// Version for test.
	Version = env.Version("1.0.0")

	// Environment for test.
	Environment = env.Environment("dev")
)

// UserAgent for test.
var UserAgent = env.UserAgent(Name.String() + "/" + Version.String())
