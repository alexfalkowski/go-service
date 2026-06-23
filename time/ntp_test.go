package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
)

func TestValidNTP(t *testing.T) {
	requireAnyNetworkNow(t,
		&time.Config{Kind: "ntp", Address: "0.beevik-ntp.pool.ntp.org", Timeout: 2 * time.Second},
		&time.Config{Kind: "ntp", Address: "1.beevik-ntp.pool.ntp.org", Timeout: 2 * time.Second},
		&time.Config{Kind: "ntp", Address: "time.cloudflare.com", Timeout: 2 * time.Second},
	)
}

func TestInvalidNTP(t *testing.T) {
	requireNetworkNowError(t, &time.Config{Kind: "ntp"})
}
