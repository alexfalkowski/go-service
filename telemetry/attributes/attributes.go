package attributes

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
)

// Key is an alias of attribute.Key.
type Key = attribute.Key

// Float64 for metrics.
func Float64(key string, value float64) attribute.KeyValue {
	return attribute.Key(key).Float64(value)
}

// Int64 for metrics.
func Int64(key string, value int64) attribute.KeyValue {
	return attribute.Key(key).Int64(value)
}

// String for metrics.
func String(key, value string) attribute.KeyValue {
	return attribute.Key(key).String(value)
}

// DBSystem for metrics.
func DBSystem(name string) attribute.KeyValue {
	return semconv.DBSystemNameKey.String(name)
}
