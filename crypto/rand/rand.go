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
func GenerateBytes(size uint32) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)

	return bytes, err
}

// GenerateString will generate using letters and symbols.
func GenerateString(size uint32) (string, error) {
	return generateString(size, letters+symbols)
}

// GenerateLetters will generate using letters.
func GenerateLetters(size uint32) (string, error) {
	return generateString(size, letters)
}

func generateString(size uint32, values string) (string, error) {
	bytes := make([]byte, size)

	for i := range size {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(values))))
		if err != nil {
			return "", err
		}

		bytes[i] = values[num.Int64()]
	}

	return string(bytes), nil
}
