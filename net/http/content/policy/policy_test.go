package policy_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http/content/policy"
	"github.com/stretchr/testify/require"
)

func TestCanDecodeMatchesSupportedContentKind(t *testing.T) {
	tests := []struct {
		kind string
		want bool
	}{
		{kind: "json", want: true},
		{kind: "yaml", want: true},
		{kind: "toml", want: false},
		{kind: "", want: true},
		{kind: "gob", want: false},
		{kind: "msgpack", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			require.Equal(t, tt.want, policy.CanDecode(tt.kind))
		})
	}
}
