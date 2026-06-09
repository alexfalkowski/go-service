package attributes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/stretchr/testify/require"
)

func TestDeploymentEnvironmentName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "prod alias", value: "prod", expected: "production"},
		{name: "production", value: "production", expected: "production"},
		{name: "stage alias", value: "stage", expected: "staging"},
		{name: "staging", value: "staging", expected: "staging"},
		{name: "qa alias", value: "qa", expected: "test"},
		{name: "test", value: "test", expected: "test"},
		{name: "testing alias", value: "testing", expected: "test"},
		{name: "dev alias", value: "dev", expected: "development"},
		{name: "development", value: "development", expected: "development"},
		{name: "unknown defaults to development", value: "local", expected: "development"},
		{name: "empty defaults to development", value: "", expected: "development"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			attribute := attributes.DeploymentEnvironmentName(tt.value)

			require.Equal(t, "deployment.environment.name", string(attribute.Key))
			require.Equal(t, tt.expected, attribute.Value.AsString())
		})
	}
}
