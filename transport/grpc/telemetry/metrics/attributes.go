package metrics

import "go.opentelemetry.io/otel/attribute"

const (
	kindAttribute   = attribute.Key("kind")
	pathAttribute   = attribute.Key("path")
	methodAttribute = attribute.Key("method")
	codeAttribute   = attribute.Key("code")
)
