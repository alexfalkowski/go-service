package rand

import (
	"crypto/rand"
	"math/big"
)

// Code is adapted from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb.

const (
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols = "~!@#$%^&*()_+-={}|[]<>?,./"
)

// GenerateBytes for rand.
func GenerateBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	return b, err
}

// GenerateString will generate using letters and symbols.
func GenerateString(n uint32) (string, error) {
	return generateString(n, letters+symbols)
}

// GenerateLetters will generate using letters.
func GenerateLetters(n uint32) (string, error) {
	return generateString(n, letters)
}

func generateString(n uint32, s string) (string, error) {
	r := make([]byte, n)

	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
		if err != nil {
			return "", err
		}

		r[i] = s[num.Int64()]
	}

	return string(r), nil
}
