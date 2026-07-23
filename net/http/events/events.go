package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol"
	protocolhttp "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// TextPlain is the CloudEvents text/plain data content type.
const TextPlain = cloudevents.TextPlain

// Event is a CloudEvent.
type Event = cloudevents.Event

// Result is a CloudEvents protocol result.
type Result = protocol.Result

// ReceiverFunc is invoked for each received CloudEvent.
type ReceiverFunc func(context.Context, Event) Result

// Client sends CloudEvents.
type Client = client.Client

// ContextWithRetriesConstantBackoff aliases the CloudEvents retry-context helper.
var ContextWithRetriesConstantBackoff = cloudevents.ContextWithRetriesConstantBackoff

// NewEvent constructs an empty CloudEvent.
func NewEvent() Event {
	return cloudevents.NewEvent()
}

// NewHTTPResult constructs a CloudEvents HTTP result with status code and message.
func NewHTTPResult(statusCode int, messageFmt string, args ...any) Result {
	return protocolhttp.NewResult(statusCode, messageFmt, args...)
}

// NewClient constructs a CloudEvents HTTP client that uses httpClient.
//
// This wrapper supports the repository's standard HTTP sender path: it supplies
// only the provided HTTP client option and does not expose CloudEvents
// constructor errors to callers.
func NewClient(httpClient http.Client) Client {
	sender, _ := cloudevents.NewClientHTTP(protocolhttp.WithClient(httpClient))
	return sender
}

// NewReceiveHandler constructs an HTTP receive handler for CloudEvents.
//
// This wrapper supports the repository's standard receiver path: it uses the
// default CloudEvents HTTP protocol with a typed ReceiverFunc and does not
// expose CloudEvents constructor errors to callers.
func NewReceiveHandler(ctx context.Context, receiver ReceiverFunc) http.Handler {
	protocol, _ := cloudevents.NewHTTP()

	handler, _ := cloudevents.NewHTTPReceiveHandler(ctx, protocol, receiver)
	return handler
}

// SendStructured sends event using structured CloudEvents encoding.
func SendStructured(ctx context.Context, client Client, event Event) Result {
	return client.Send(binding.WithForceStructured(ctx), event)
}

// SendBinary sends event using binary CloudEvents encoding.
func SendBinary(ctx context.Context, client Client, event Event) Result {
	return client.Send(binding.WithForceBinary(ctx), event)
}

// ContextWithTarget returns a context with the CloudEvents target set.
func ContextWithTarget(ctx context.Context, target string) context.Context {
	return cloudevents.ContextWithTarget(ctx, target)
}

// IsACK reports whether result is an ACK.
func IsACK(result Result) bool {
	return protocol.IsACK(result)
}

// IsNACK reports whether result is a NACK.
func IsNACK(result Result) bool {
	return protocol.IsNACK(result)
}
