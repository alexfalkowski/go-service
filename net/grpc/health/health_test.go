package health_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/stretchr/testify/require"
)

func TestIsMethodName(t *testing.T) {
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
			require.Equal(t, tt.match, health.IsMethodName(tt.method))
		})
	}
}
