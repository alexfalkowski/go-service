package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// SchemaURL is the OpenTelemetry semantic conventions schema URL used by this package.
//
// It is an alias of semconv.SchemaURL and is typically passed to
// resource.NewWithAttributes to indicate which schema version the provided
// attributes follow.
const SchemaURL = semconv.SchemaURL

// Key is an alias of attribute.Key.
type Key = attribute.Key

// KeyValue is an alias of attribute.KeyValue.
type KeyValue = attribute.KeyValue

// RPCSystemNameGRPC identifies gRPC as the RPC system for semantic conventions.
//
// It is an alias of semconv.RPCSystemNameGRPC and is commonly used by
// instrumentation to label RPC-related telemetry consistently.
var RPCSystemNameGRPC = semconv.RPCSystemNameGRPC

// HostID returns a host.id attribute with value val.
//
// This is a thin wrapper around semconv.HostID and is typically used when
// constructing a Resource to describe the running service instance.
func HostID(val string) attribute.KeyValue {
	return semconv.HostID(val)
}

// ServiceName returns a service.name attribute with value val.
//
// This is a thin wrapper around semconv.ServiceName and is typically used when
// constructing a Resource to describe the running service.
func ServiceName(val string) attribute.KeyValue {
	return semconv.ServiceName(val)
}

// ServiceVersion returns a service.version attribute with value val.
//
// This is a thin wrapper around semconv.ServiceVersion and is typically used when
// constructing a Resource to describe the running service.
func ServiceVersion(val string) attribute.KeyValue {
	return semconv.ServiceVersion(val)
}

// DeploymentEnvironmentName returns a deployment.environment.name attribute with value val.
//
// This is a thin wrapper around semconv.DeploymentEnvironmentName and is typically
// used when constructing a Resource to describe the deployment environment (for
// example "prod", "staging").
func DeploymentEnvironmentName(val string) attribute.KeyValue {
	return semconv.DeploymentEnvironmentName(val)
}
