package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
)

// Generate key pair with RSA.
func Generate() (PublicKey, PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&p.PublicKey)})
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(p)})

	return PublicKey(pub), PrivateKey(pri), nil
}

// Algo for rsa.
type Algo interface {
	algo.Cipher
}

// NewAlgo for rsa.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &algo.NoCipher{}, nil
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		return nil, err
	}

	pri, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}

	return &rsaAlgo{publicKey: pub, privateKey: pri}, nil
}

type rsaAlgo struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (a *rsaAlgo) Encrypt(msg string) (string, error) {
	e, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, a.publicKey, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(e), err
}

func (a *rsaAlgo) Decrypt(msg string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	d, err = rsa.DecryptOAEP(sha512.New(), rand.Reader, a.privateKey, d, nil)

	return string(d), err
}
