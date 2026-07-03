package attributes_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/stretchr/testify/require"
)

func TestResourceMergesConfiguredAttributesWithIdentity(t *testing.T) {
	t.Parallel()

	resource := attributes.NewResource(
		attributes.Map{
			"k8s.namespace.name":                                 "payments",
			string(attributes.HostID("").Key):                    "configured",
			string(attributes.ServiceName("").Key):               "configured",
			string(attributes.ServiceVersion("").Key):            "configured",
			string(attributes.DeploymentEnvironmentName("").Key): "configured",
		},
		"host-id",
		"service-name",
		"service-version",
		"prod",
	)
	attrs := resourceAttributes(resource.Attributes())

	require.Equal(t, "payments", attrs["k8s.namespace.name"])
	require.Equal(t, "host-id", attrs["host.id"])
	require.Equal(t, "service-name", attrs["service.name"])
	require.Equal(t, "service-version", attrs["service.version"])
	require.Equal(t, "production", attrs["deployment.environment.name"])
}

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

func resourceAttributes(attrs []attributes.KeyValue) map[string]string {
	values := make(map[string]string)
	for _, attr := range attrs {
		values[string(attr.Key)] = attr.Value.AsString()
	}
	return values
}
