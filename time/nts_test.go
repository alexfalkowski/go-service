package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
)

func TestValidNTS(t *testing.T) {
	requireAnyNetworkNow(t,
		&time.Config{Kind: "nts", Address: "time.cloudflare.com", Timeout: 2 * time.Second},
		&time.Config{Kind: "nts", Address: "nts.netnod.se", Timeout: 2 * time.Second},
	)
}

func TestInvalidNTS(t *testing.T) {
	requireNetworkNowError(t, &time.Config{Kind: "nts"})
}
