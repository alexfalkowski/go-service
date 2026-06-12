// Package events provides CloudEvents HTTP helpers behind the go-service import path.
//
// This package owns the repository's direct dependency on the CloudEvents SDK for HTTP event
// sender/receiver wiring. Higher-level transport packages should import this package rather than
// importing github.com/cloudevents/sdk-go/v2 directly.
package events
