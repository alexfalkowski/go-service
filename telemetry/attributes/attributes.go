package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

// SchemaURL is the OpenTelemetry semantic conventions schema URL used by this package.
//
// It is an alias of semconv.SchemaURL and is typically passed to
// resource.NewWithAttributes to declare which semantic-convention schema version
// the provided attributes follow.
const SchemaURL = semconv.SchemaURL

// Key aliases attribute.Key for callers that want to work with OpenTelemetry
// attribute keys through this package without importing attribute directly.
//
// It is primarily useful when building additional semantic-convention helpers in
// packages that already depend on telemetry/attributes.
type Key = attribute.Key

// KeyValue aliases attribute.KeyValue for callers that want to exchange
// OpenTelemetry attributes through this package without importing attribute
// directly.
//
// Functions in this package return KeyValue-compatible values.
type KeyValue = attribute.KeyValue

// DBSystemNamePostgreSQL identifies PostgreSQL as the value of the
// db.system.name semantic-convention attribute.
//
// It is an alias of semconv.DBSystemNamePostgreSQL and is a fully-formed
// attribute.KeyValue that can be attached directly to spans, metrics, logs, or
// resources.
var DBSystemNamePostgreSQL = semconv.DBSystemNamePostgreSQL

// RPCSystemNameGRPC identifies gRPC as the RPC system for semantic conventions.
//
// It is an alias of semconv.RPCSystemNameGRPC and is commonly used by
// instrumentation to label RPC-related telemetry consistently. The exported
// value is already an attribute.KeyValue and can be attached directly.
var RPCSystemNameGRPC = semconv.RPCSystemNameGRPC

// DBSystemNameKey identifies the semantic convention key for a database system.
//
// It is an alias of semconv.DBSystemNameKey and is commonly used by SQL
// instrumentation to attach the database system name consistently.
//
// Callers typically use DBSystemNameKey.String("postgresql") or reuse one of the
// predeclared attribute values such as DBSystemNamePostgreSQL.
var DBSystemNameKey = semconv.DBSystemNameKey

// HostID returns a host.id attribute with value val.
//
// This is a thin wrapper around semconv.HostID and is typically used when
// constructing a Resource to describe the running service instance.
//
// Parameters:
//   - val: the host identifier reported for the current process or machine
func HostID(val string) attribute.KeyValue {
	return semconv.HostID(val)
}

// ServiceName returns a service.name attribute with value val.
//
// This is a thin wrapper around semconv.ServiceName and is typically used when
// constructing a Resource to describe the running service.
//
// Parameters:
//   - val: the logical service name to report in telemetry
func ServiceName(val string) attribute.KeyValue {
	return semconv.ServiceName(val)
}

// ServiceVersion returns a service.version attribute with value val.
//
// This is a thin wrapper around semconv.ServiceVersion and is typically used when
// constructing a Resource to describe the running service.
//
// Parameters:
//   - val: the version string to report for the running service build
func ServiceVersion(val string) attribute.KeyValue {
	return semconv.ServiceVersion(val)
}

// DeploymentEnvironmentName returns a deployment.environment.name attribute with value val.
//
// This is a thin wrapper around semconv.DeploymentEnvironmentName and is typically
// used when constructing a Resource to describe the deployment environment (for
// example "prod", "staging").
//
// Parameters:
//   - val: the deployment environment name, such as "prod" or "staging"
func DeploymentEnvironmentName(val string) attribute.KeyValue {
	return semconv.DeploymentEnvironmentName(val)
}
