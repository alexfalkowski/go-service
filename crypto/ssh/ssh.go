package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/pem"

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
	// Encrypt msg.
	Encrypt(msg string) (string, error)

	// Decrypt msg.
	Decrypt(msg string) (string, error)
}

// NewAlgo for ssh.
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
	d, err := cfg.GetPublic()
	if err != nil {
		return nil, err
	}

	//nolint:dogsled
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(d)
	if err != nil {
		return nil, err
	}

	key := parsed.(ssh.CryptoPublicKey)

	return key.CryptoPublicKey().(*rsa.PublicKey), nil
}

func privateKey(cfg *Config) (*rsa.PrivateKey, error) {
	d, err := cfg.GetPrivate()
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(d)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}
