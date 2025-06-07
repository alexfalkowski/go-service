package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

type (
	// Key is an alias of attribute.Key.
	Key = attribute.Key

	// KeyValue is an alias of attribute.KeyValue.
	KeyValue = attribute.KeyValue
)

var (
	// HostID is an alias of semconv.HostID.
	HostID = semconv.HostID

	// HTTPResponseStatusCode is an alias of semconv.HTTPResponseStatusCode.
	HTTPResponseStatusCode = semconv.HTTPResponseStatusCode

	// HTTPRoute is an alias of semconv.HTTPRoute.
	HTTPRoute = semconv.HTTPRoute

	// RPCSystemGRPC is an alias of semconv.RPCSystemGRPC.
	RPCSystemGRPC = semconv.RPCSystemGRPC

	// RPCService is an alias of semconv.RPCService.
	RPCService = semconv.RPCService

	// RPCMethod is an alias of semconv.RPCMethod.
	RPCMethod = semconv.RPCMethod

	// SchemaURL is an alias of semconv.SchemaURL.
	SchemaURL = semconv.SchemaURL

	// ServiceName is an alias of semconv.ServiceName.
	ServiceName = semconv.ServiceName

	// ServiceVersion is an alias of semconv.ServiceVersion.
	ServiceVersion = semconv.ServiceVersion

	// DeploymentEnvironmentName is an alias of semconv.DeploymentEnvironmentName.
	DeploymentEnvironmentName = semconv.DeploymentEnvironmentName
)

// Bool for attributes.
func Bool(key string, value bool) attribute.KeyValue {
	return attribute.Key(key).Bool(value)
}

// Float64 for attributes.
func Float64(key string, value float64) attribute.KeyValue {
	return attribute.Key(key).Float64(value)
}

// Int64 for attributes.
func Int64(key string, value int64) attribute.KeyValue {
	return attribute.Key(key).Int64(value)
}

// String for attributes.
func String(key, value string) attribute.KeyValue {
	return attribute.Key(key).String(value)
}

// DBSystem for attributes.
func DBSystem(name string) attribute.KeyValue {
	return semconv.DBSystemNameKey.String(name)
}

// HTTPRequestMethod for attributes.
func HTTPRequestMethod(name string) attribute.KeyValue {
	return semconv.HTTPRequestMethodKey.String(name)
}

// GRPCStatusCode for attributes.
func GRPCStatusCode(code int64) attribute.KeyValue {
	return semconv.RPCGRPCStatusCodeKey.Int64(code)
}
