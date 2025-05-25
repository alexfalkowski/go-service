package test

import "github.com/alexfalkowski/go-service/v2/crypto/ssh"

// NewSSH for test.
func NewSSH(public, private string) *ssh.Config {
	return &ssh.Config{
		Public:  FilePath(public),
		Private: FilePath(private),
	}
}
