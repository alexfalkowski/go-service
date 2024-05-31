package ssh

import (
	"net"
	"time"

	"github.com/alexfalkowski/go-service/errors"
	t "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/ssh/handler"
	logger "github.com/alexfalkowski/go-service/transport/ssh/telemetry/logger/zap"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// ClientOption for SSH.
type ClientOption interface{ apply(opts *clientOpts) }

type clientOpts struct {
	logger  *zap.Logger
	timeout time.Duration
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) { f(o) }

// WithClientTimeout for SSH.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = t.MustParseDuration(timeout)
	})
}

// WithClientLogger for SSH.
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// NewClient for SSH.
func NewClient(target string, opts ...ClientOption) (*Client, error) {
	os := clientOptions(opts...)

	cfg := &ssh.ClientConfig{
		Timeout:         os.timeout,
		HostKeyCallback: hostKey,
	}

	c, err := ssh.Dial("tcp", target, cfg)
	if err != nil {
		return nil, errors.Prefix("ssh client", err)
	}

	var client handler.Client = &client{client: c}

	if os.logger != nil {
		client = logger.NewClient(os.logger, client)
	}

	return &Client{ssh: c, client: client}, nil
}

// Client fot SSH.
type Client struct {
	ssh    *ssh.Client
	client handler.Client
}

// Run the command.
func (c *Client) Run(cmd string) ([]byte, error) {
	b, err := c.client.Run(cmd)

	return b, errors.Prefix("ssh run", err)
}

// Close the client.
func (c *Client) Close() error {
	return errors.Prefix("ssh close", c.ssh.Close())
}

type client struct {
	client *ssh.Client
}

func (c *client) Run(cmd string) ([]byte, error) {
	s, err := c.client.NewSession()
	if err != nil {
		return nil, errors.Prefix("new ssh session", err)
	}

	return s.Output(cmd)
}

func clientOptions(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	return os
}

func hostKey(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}
