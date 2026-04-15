package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Connect retries dialing address for up to one second until it succeeds.
//
// The address may be a raw host:port or a go-service "<network>://<address>" string.
func Connect(ctx context.Context, address string) (net.Conn, error) {
	network, addr := net.ListenNetworkAddress(address)
	dialer := &net.Dialer{}
	deadline := time.Now().Add(time.Second.Duration())
	var err error

	for time.Now().Before(deadline) {
		conn, dialErr := dialer.DialContext(ctx, network, addr)
		if dialErr == nil {
			return conn, nil
		}

		err = dialErr
		time.Sleep(10 * time.Millisecond)
	}

	return nil, err
}
