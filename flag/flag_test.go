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

func TestAddConfigParsesLongAndShortConfigFlags(t *testing.T) {
	cases := []struct {
		name string
		want string
		args []string
	}{
		{name: "default", want: "file:config.yaml"},
		{name: "long flag", args: []string{"-config", "file:override.yaml"}, want: "file:override.yaml"},
		{name: "short flag", args: []string{"-c", "env:CONFIG"}, want: "env:CONFIG"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet("test")
			set.AddConfig("file:config.yaml")

			require.NoError(t, set.Parse(tt.args))
			require.Equal(t, tt.want, set.GetConfig())
		})
	}
}
