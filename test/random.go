package test

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	min = 2000
	max = 65535
)

// GenerateRandomPort for test.
func GenerateRandomPort() string {
	rand.Seed(time.Now().UnixNano())

	port := rand.Intn(max-min+1) + min // nolint:gosec

	return strconv.Itoa(port)
}
