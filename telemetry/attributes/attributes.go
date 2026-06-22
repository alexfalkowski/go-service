package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
)

// SchemaURL is the OpenTelemetry semantic conventions schema URL used by this package.
//
// It is an alias of [semconv.SchemaURL] and is typically passed to
// [go.opentelemetry.io/otel/sdk/resource.NewWithAttributes] to declare which semantic-convention schema version
// the provided attributes follow.
const SchemaURL = semconv.SchemaURL

// Key aliases [attribute.Key] for callers that want to work with OpenTelemetry
// attribute keys through this package without importing attribute directly.
//
// It is primarily useful when building additional semantic-convention helpers in
// packages that already depend on telemetry/attributes.
type Key = attribute.Key

// KeyValue aliases [attribute.KeyValue] for callers that want to exchange
// OpenTelemetry attributes through this package without importing attribute
// directly.
//
// Functions in this package return KeyValue-compatible values.
type KeyValue = attribute.KeyValue

// DBSystemNamePostgreSQL identifies PostgreSQL as the value of the
// db.system.name semantic-convention attribute.
//
// It is an alias of [semconv.DBSystemNamePostgreSQL] and is a fully-formed
// [attribute.KeyValue] that can be attached directly to spans, metrics, logs, or
// resources.
var DBSystemNamePostgreSQL = semconv.DBSystemNamePostgreSQL

// RPCSystemNameGRPC identifies gRPC as the RPC system for semantic conventions.
//
// It is an alias of [semconv.RPCSystemNameGRPC] and is commonly used by
// instrumentation to label RPC-related telemetry consistently. The exported
// value is already an [attribute.KeyValue] and can be attached directly.
var RPCSystemNameGRPC = semconv.RPCSystemNameGRPC

// DBSystemNameKey identifies the semantic convention key for a database system.
//
// It is an alias of [semconv.DBSystemNameKey] and is commonly used by SQL
// instrumentation to attach the database system name consistently.
//
// Callers typically use [DBSystemNameKey.String]("postgresql") or reuse one of the
// predeclared attribute values such as DBSystemNamePostgreSQL.
var DBSystemNameKey = semconv.DBSystemNameKey

// DBClientConnectionPoolName returns a db.client.connection.pool.name attribute
// with value name.
//
// This is a thin wrapper around [semconv.DBClientConnectionPoolName] and is
// typically used when registering SQL connection-pool metrics.
//
// Parameters:
//   - name: the connection pool name, unique within the instrumented
//     application
func DBClientConnectionPoolName(name string) attribute.KeyValue {
	return semconv.DBClientConnectionPoolName(name)
}

// HostID returns a host.id attribute with value val.
//
// This is a thin wrapper around [semconv.HostID] and is typically used when
// constructing a Resource to describe the running service instance.
//
// Parameters:
//   - val: the host identifier reported for the current process or machine
func HostID(val string) attribute.KeyValue {
	return semconv.HostID(val)
}

// ServiceName returns a service.name attribute with value name.
//
// This is a thin wrapper around [semconv.ServiceName] and is typically used when
// constructing a Resource to describe the running service.
//
// Parameters:
//   - name: the logical service name to report in telemetry
func ServiceName(name string) attribute.KeyValue {
	return semconv.ServiceName(name)
}

// ServiceVersion returns a service.version attribute with value version.
//
// This is a thin wrapper around [semconv.ServiceVersion] and is typically used when
// constructing a Resource to describe the running service.
//
// Parameters:
//   - version: the version string to report for the running service build
func ServiceVersion(version string) attribute.KeyValue {
	return semconv.ServiceVersion(version)
}

// DeploymentEnvironmentName returns a deployment.environment.name attribute for env.
//
// It maps common environment names onto the stable OpenTelemetry enum values:
// "prod" and "production" become "production"; "stage" and "staging" become
// "staging"; "qa", "test", and "testing" become "test"; "dev" and
// "development" become "development". Unknown or empty values default to
// "development".
//
// It is typically used when constructing a Resource to describe the deployment
// environment.
//
// Parameters:
//   - env: the deployment environment name, such as "prod" or "staging"
func DeploymentEnvironmentName(env string) attribute.KeyValue {
	switch env {
	case "prod", "production":
		return semconv.DeploymentEnvironmentNameProduction
	case "stage", "staging":
		return semconv.DeploymentEnvironmentNameStaging
	case "qa", "test", "testing":
		return semconv.DeploymentEnvironmentNameTest
	case "dev", "development":
		return semconv.DeploymentEnvironmentNameDevelopment
	}

	return semconv.DeploymentEnvironmentNameDevelopment
}
