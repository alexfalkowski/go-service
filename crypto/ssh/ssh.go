package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"golang.org/x/crypto/ssh"
)

// Generate key pair with ssh.
func Generate() (PublicKey, PrivateKey, error) {
	p, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	ppr, err := ssh.MarshalPrivateKey(p, "")
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(p.Public())
	if err != nil {
		return "", "", err
	}

	return PublicKey(ssh.MarshalAuthorizedKey(pub)), PrivateKey(pem.EncodeToMemory(ppr)), nil
}

// Algo for ssh.
type Algo interface {
	algo.Cipher
}

// NewAlgo for ssh.
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

	return &sshAlgo{publicKey: pub, privateKey: pri}, nil
}

type sshAlgo struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func (a *sshAlgo) Encrypt(msg string) (string, error) {
	e, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, a.publicKey, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(e), err
}

func (a *sshAlgo) Decrypt(msg string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	d, err = rsa.DecryptOAEP(sha512.New(), rand.Reader, a.privateKey, d, nil)

	return string(d), err
}
