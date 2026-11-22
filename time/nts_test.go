package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/time"
)

func TestNetworkNowReturnsTimeFromReachableNTSServer(t *testing.T) {
	requireAnyNetworkNow(t,
		&time.Config{Kind: "nts", Endpoint: "time.cloudflare.com", Timeout: 2 * time.Second},
		&time.Config{Kind: "nts", Endpoint: "nts.netnod.se", Timeout: 2 * time.Second},
	)
}

func TestNetworkNowReturnsErrorForIncompleteNTSConfig(t *testing.T) {
	requireNetworkNowError(t, &time.Config{Kind: "nts"})
}
