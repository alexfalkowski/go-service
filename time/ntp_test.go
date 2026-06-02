package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
)

func TestValidNTP(t *testing.T) {
	requireNetworkNow(t, &time.Config{Kind: "ntp", Address: "0.beevik-ntp.pool.ntp.org"})
}

func TestInvalidNTP(t *testing.T) {
	requireNetworkNowError(t, &time.Config{Kind: "ntp"})
}
