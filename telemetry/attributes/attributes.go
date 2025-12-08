package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

// SchemaURL is an alias of semconv.SchemaURL.
const SchemaURL = semconv.SchemaURL

type (
	// Key is an alias of attribute.Key.
	Key = attribute.Key

	// KeyValue is an alias of attribute.KeyValue.
	KeyValue = attribute.KeyValue
)

// RPCSystemGRPC is an alias of semconv.RPCSystemGRPC.
var RPCSystemGRPC = semconv.RPCSystemGRPC

// HostID is an alias of semconv.HostID.
func HostID(val string) attribute.KeyValue {
	return semconv.HostID(val)
}

// ServiceName is an alias of semconv.ServiceName.
func ServiceName(val string) attribute.KeyValue {
	return semconv.ServiceName(val)
}

// ServiceVersion is an alias of semconv.ServiceVersion.
func ServiceVersion(val string) attribute.KeyValue {
	return semconv.ServiceVersion(val)
}

// DeploymentEnvironmentName is an alias of semconv.DeploymentEnvironmentName.
func DeploymentEnvironmentName(val string) attribute.KeyValue {
	return semconv.DeploymentEnvironmentName(val)
}
