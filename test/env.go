package test

import (
	"github.com/alexfalkowski/go-service/env"
)

var (
	// DevEnvironment for test.
	DevEnvironment = env.Environment("dev")

	// ProdEnvironment for test.
	ProdEnvironment = env.Environment("prod")
)
