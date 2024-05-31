package handler

import (
	"github.com/gliderlabs/ssh"
)

// Server for SSH.
type Server interface {
	// Handle command.
	Handle(ctx ssh.Context, cmd []string) error
}

// Client for SSH.
type Client interface {
	// Run command.
	Run(cmd string) ([]byte, error)
}
