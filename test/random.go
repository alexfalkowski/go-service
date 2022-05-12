package test

import (
	"math/rand"
	"strconv"
)

const (
	min = 10000
	max = 65535
)

// GenerateRandomPort for test.
func GenerateRandomPort() string {
	port := rand.Intn(max-min+1) + min // #nosec G404

	return strconv.Itoa(port)
}
