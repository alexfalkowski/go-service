package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
)

// Code was adapted from github.com/alexellis/hmac/v2.

// Generate for hmac.
func Generate() (string, error) {
	s, err := rand.GenerateBytes(32)

	return base64.StdEncoding.EncodeToString(s), err
}

// Algo for hmac.
type Algo interface {
	algo.Signer
}

// NewAlgo for hmac.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &algo.NoSigner{}, nil
	}

	k, err := cfg.GetKey()

	return &hmacAlgo{key: []byte(k)}, err
}

type hmacAlgo struct {
	key []byte
}

func (a *hmacAlgo) Sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, a.key)

	_, err := mac.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (a *hmacAlgo) Verify(sig, msg string) error {
	decoded, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	mac := hmac.New(sha512.New, a.key)

	_, err = mac.Write([]byte(msg))
	if err != nil {
		return err
	}

	expectedMAC := mac.Sum(nil)

	if !hmac.Equal(decoded, expectedMAC) {
		return errors.ErrInvalidMatch
	}

	return nil
}
