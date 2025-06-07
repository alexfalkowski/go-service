package attributes

import "go.opentelemetry.io/otel/attribute"

// Key is an alias of attribute.Key.
type Key = attribute.Key

// String for metrics.
func String(key, value string) attribute.KeyValue {
	return attribute.Key(key).String(value)
}
