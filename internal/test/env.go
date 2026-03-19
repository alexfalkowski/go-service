package test

import (
	"github.com/alexfalkowski/go-service/v2/env"
)

const (
	// ID is the shared deterministic ID fixture used in tests.
	ID = env.ID("1234567890")

	// Name is the shared service name used by test helpers and route registration.
	Name = env.Name("test")

	// UserID is the shared user identifier derived from Name.
	UserID = env.UserID(Name)

	// Version is the shared service version used by test helpers.
	Version = env.Version("1.0.0")

	// Environment is the shared deployment environment label used by test telemetry.
	Environment = env.Environment("dev")
)

// UserAgent is the shared user agent string derived from Name and Version.
var UserAgent = env.UserAgent(Name.String() + "/" + Version.String())
