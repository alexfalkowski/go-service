package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/stretchr/testify/require"
)

func TestIsOperationMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		match  bool
	}{
		{name: "grpc health check", method: health.CheckFullMethodName, match: true},
		{name: "grpc health watch", method: health.WatchFullMethodName, match: true},
		{name: "grpc health list", method: health.ListFullMethodName, match: true},
		{name: "custom service with health in name", method: "/customer.health.v1.HealthPlans/Check", match: false},
		{name: "custom method with metrics in name", method: "/customer.v1.Report/GetMetrics", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, strings.IsOperationMethod(tt.method))
		})
	}
}

func TestSplitServiceMethod(t *testing.T) {
	tests := []struct {
		name    string
		full    string
		service string
		method  string
		ok      bool
	}{
		{name: "full method", full: "/greet.v1.Greeter/SayHello", service: "greet.v1.Greeter", method: "SayHello", ok: true},
		{name: "missing package", full: "/Greeter/SayHello"},
		{name: "missing leading slash", full: "greet.v1.Greeter/SayHello"},
		{name: "missing method", full: "/greet.v1.Greeter"},
		{name: "extra path segment", full: "/greet.v1.Greeter/SayHello/Again"},
		{name: "empty method", full: "/greet.v1.Greeter/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, method, ok := strings.SplitServiceMethod(test.full)
			require.Equal(t, test.service, service)
			require.Equal(t, test.method, method)
			require.Equal(t, test.ok, ok)
		})
	}
}
