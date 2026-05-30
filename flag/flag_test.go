package flag_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/stretchr/testify/require"
)

func TestGetConfigWithoutAddConfig(t *testing.T) {
	set := flag.NewFlagSet("test")

	require.Empty(t, set.GetConfig())
}

func TestAddConfig(t *testing.T) {
	tests := []struct {
		name string
		want string
		args []string
	}{
		{name: "default", want: "file:config.yml"},
		{name: "long flag", args: []string{"-config", "file:override.yml"}, want: "file:override.yml"},
		{name: "short flag", args: []string{"-c", "env:CONFIG"}, want: "env:CONFIG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet("test")
			set.AddConfig("file:config.yml")

			require.NoError(t, set.Parse(tt.args))
			require.Equal(t, tt.want, set.GetConfig())
		})
	}
}
