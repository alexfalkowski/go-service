package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/stretchr/testify/require"
)

func FuzzParseServiceMethod(f *testing.F) {
	for _, full := range []string{
		"/greet.v1.Greeter/SayHello",
		"greet.v1.Greeter/SayHello",
		"/greet.v1.Greeter",
		"/greet.v1.Greeter/",
		"//SayHello",
		"",
	} {
		f.Add(full)
	}

	f.Fuzz(func(t *testing.T, full string) {
		service, method := grpc.ParseServiceMethod(full)
		splitService, splitMethod, ok := url.SplitPath(full)
		if !ok {
			require.Equal(t, "root", service)
			require.Equal(t, "root", method)
			return
		}

		require.Equal(t, splitService, service)
		require.Equal(t, splitMethod, method)
	})
}
