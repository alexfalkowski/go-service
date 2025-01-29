package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	kindAttribute    = attribute.Key("kind")
	serviceAttribute = attribute.Key("service")
	methodAttribute  = attribute.Key("method")
	codeAttribute    = attribute.Key("code")
)
