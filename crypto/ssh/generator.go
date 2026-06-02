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

// NewGenerator constructs a Generator that produces Ed25519 SSH key pairs.
//
// The provided generator is used as the cryptographically-secure randomness source for key generation.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates Ed25519 SSH key pairs.
type Generator struct {
	generator *rand.Generator
}

// Generate returns an Ed25519 public/private key pair encoded in SSH formats.
//
// The returned values are compatible with the expectations of [Config]:
//
//   - public: SSH authorized_keys format (via x/crypto/ssh.MarshalAuthorizedKey)
//   - private: PEM-encoded SSH private key (via x/crypto/ssh.MarshalPrivateKey and pem.EncodeToMemory)
//
// If key generation, marshaling, or encoding fails, the returned error is prefixed with "ssh".
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.generator)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	publicKey, err := ssh.NewPublicKey(public)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	privateBlock, err := ssh.MarshalPrivateKey(private, strings.Empty)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	pub := ssh.MarshalAuthorizedKey(publicKey)
	pri := pem.EncodeToMemory(privateBlock)
	return bytes.String(pub), bytes.String(pri), nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("ssh", err)
}
