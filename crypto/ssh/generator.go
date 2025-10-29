package ssh

import (
	"crypto/ed25519"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
	"golang.org/x/crypto/ssh"
)

// NewGenerator for ssh.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator for ssh.
type Generator struct {
	generator *rand.Generator
}

// Generate key pair with ssh.
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.generator)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	mpu, err := ssh.NewPublicKey(public)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	mpr, err := ssh.MarshalPrivateKey(private, "")
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	pub := ssh.MarshalAuthorizedKey(mpu)
	pri := pem.EncodeToMemory(mpr)
	return bytes.String(pub), bytes.String(pri), nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("ssh", err)
}
