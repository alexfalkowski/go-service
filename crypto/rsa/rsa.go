package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
)

// Generate key pair with RSA.
func Generate() (PublicKey, PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	pub := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&p.PublicKey))
	pri := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(p))

	return PublicKey(pub), PrivateKey(pri), nil
}

// Algo for rsa.
type Algo interface {
	// Encrypt msg.
	Encrypt(msg string) (string, error)

	// Decrypt msg.
	Decrypt(msg string) (string, error)
}

// NewAlgo for rsa.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &none{}, nil
	}

	pub, err := publicKey(cfg)
	if err != nil {
		return nil, err
	}

	pri, err := privateKey(cfg)
	if err != nil {
		return nil, err
	}

	return &algo{publicKey: pub, privateKey: pri}, nil
}

type algo struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (a *algo) Encrypt(msg string) (string, error) {
	e, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, a.publicKey, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(e), err
}

func (a *algo) Decrypt(msg string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	d, err = rsa.DecryptOAEP(sha512.New(), rand.Reader, a.privateKey, d, nil)

	return string(d), err
}

type none struct{}

func (*none) Encrypt(msg string) (string, error) {
	return msg, nil
}

func (*none) Decrypt(msg string) (string, error) {
	return msg, nil
}

func publicKey(cfg *Config) (*rsa.PublicKey, error) {
	k, err := cfg.GetPublic()
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey([]byte(k))
}

func privateKey(cfg *Config) (*rsa.PrivateKey, error) {
	k, err := cfg.GetPrivate()
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey([]byte(k))
}
