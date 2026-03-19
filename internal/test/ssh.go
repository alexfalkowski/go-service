package test

import "github.com/alexfalkowski/go-service/v2/crypto/ssh"

// NewSSH returns an SSH key config that resolves the supplied public and private key fixtures from `test/`.
func NewSSH(public, private string) *ssh.Config {
	return &ssh.Config{
		Public:  FilePath(public),
		Private: FilePath(private),
	}
}
