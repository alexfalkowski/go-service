package ssh

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
	cerr "github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/runtime"
	"golang.org/x/crypto/ssh"
)

// NewGenerator for ssh.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for ssh.
type Generator struct {
	gen *rand.Generator
}

// Generate key pair with ssh.
func (g *Generator) Generate() (pub string, pri string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("ssh", runtime.ConvertRecover(r))
		}
	}()

	public, private, err := ed25519.GenerateKey(g.gen)
	runtime.Must(err)

	mpu, err := ssh.NewPublicKey(public)
	runtime.Must(err)

	mpr, err := ssh.MarshalPrivateKey(private, "")
	runtime.Must(err)

	pub = string(ssh.MarshalAuthorizedKey(mpu))
	pri = string(pem.EncodeToMemory(mpr))

	return
}

// NewSigner for ssh.
func NewSigner(cfg *Config) (Signer, error) {
	if !IsEnabled(cfg) {
		return &algo.NoSigner{}, nil
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		return nil, err
	}

	pri, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}

	return &signer{publicKey: pub, privateKey: pri}, nil
}

type (
	// Signer for ssh.
	Signer interface {
		algo.Signer
	}

	signer struct {
		publicKey  ed25519.PublicKey
		privateKey ed25519.PrivateKey
	}
)

func (a *signer) Sign(msg string) (string, error) {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m), nil
}

func (a *signer) Verify(sig, msg string) error {
	d, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.publicKey, []byte(msg), d)
	if !ok {
		return cerr.ErrInvalidMatch
	}

	return nil
}
