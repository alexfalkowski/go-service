package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

var (
	serviceAttribute    = attribute.Key("service")
	methodAttribute     = attribute.Key("method")
	statusCodeAttribute = attribute.Key("status_code")
)
