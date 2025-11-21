package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
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

// RPCService is an alias of semconv.RPCService.
func RPCService(val string) attribute.KeyValue {
	return semconv.RPCService(val)
}

// RPCMethod is an alias of semconv.RPCMethod.
func RPCMethod(val string) attribute.KeyValue {
	return semconv.RPCMethod(val)
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

// GRPCStatusCode for attributes.
func GRPCStatusCode(code int64) attribute.KeyValue {
	return semconv.RPCGRPCStatusCodeKey.Int64(code)
}
