package time_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/time"
)

func TestValidNTS(t *testing.T) {
	requireNetworkNow(t, &time.Config{Kind: "nts", Address: "time.cloudflare.com"})
}

func TestInvalidNTS(t *testing.T) {
	requireNetworkNowError(t, &time.Config{Kind: "nts"})
}
