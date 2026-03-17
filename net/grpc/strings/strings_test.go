package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/stretchr/testify/require"
)

func TestIsIgnorableFullMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		match  bool
	}{
		{name: "grpc health check", method: "/grpc.health.v1.Health/Check", match: true},
		{name: "grpc health watch", method: "/grpc.health.v1.Health/Watch", match: true},
		{name: "grpc health list", method: "/grpc.health.v1.Health/List", match: true},
		{name: "custom service with health in name", method: "/customer.health.v1.HealthPlans/Check", match: false},
		{name: "custom method with metrics in name", method: "/customer.v1.Report/GetMetrics", match: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.match, strings.IsIgnorable(tt.method))
		})
	}
}
