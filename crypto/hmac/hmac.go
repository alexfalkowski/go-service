package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
)

// Code was adapted from github.com/alexellis/hmac/v2.

// Generate for hmac.
func Generate() (Key, error) {
	s, err := rand.GenerateBytes(32)

	return Key(base64.StdEncoding.EncodeToString(s)), err
}

// Algo for hmac.
type Algo interface {
	// Generate an encoded msg.
	Generate(msg string) string

	// Compare encoded with msg.
	Compare(enc, msg string) error
}

// NewAlgo for hmac.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &none{}, nil
	}

	k, err := cfg.GetKey()

	return &algo{key: k}, err
}

type algo struct {
	key []byte
}

func (a *algo) Generate(msg string) string {
	mac := hmac.New(sha512.New, a.key)
	mac.Write([]byte(msg))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (a *algo) Compare(enc, msg string) error {
	d, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return err
	}

	mac := hmac.New(sha512.New, a.key)
	mac.Write([]byte(msg))

	expectedMAC := mac.Sum(nil)

	if !hmac.Equal(d, expectedMAC) {
		return errors.ErrMismatch
	}

	return nil
}

type none struct{}

func (*none) Generate(msg string) string {
	return msg
}

func (*none) Compare(_, _ string) error {
	return nil
}
