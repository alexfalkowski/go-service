package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	serviceAttribute    = attribute.Key("service")
	methodAttribute     = attribute.Key("method")
	statusCodeAttribute = attribute.Key("status_code")
)
