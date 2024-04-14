package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

var (
	kindAttribute    = attribute.Key("kind")
	serviceAttribute = attribute.Key("service")
	methodAttribute  = attribute.Key("method")
	codeAttribute    = attribute.Key("code")
)
