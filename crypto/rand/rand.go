package rand

import (
	"crypto/rand"
	"math/big"
)

// Code is adapted from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb.

const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols = "~!@#$%^&*()_+-={}|[]<>?,./"
	all     = letters + symbols
)

// GenerateBytes for rand.
func GenerateBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	return b, err
}

// GenerateString for rand.
func GenerateString(n uint32) (string, error) {
	r := make([]byte, n)

	for i := 0; i < int(n); i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}

		r[i] = all[num.Int64()]
	}

	return string(r), nil
}
