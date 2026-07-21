// Package events provides CloudEvents HTTP sender/receiver wiring for go-service.
//
// This package integrates go-service CloudEvents HTTP helpers with transport wiring.
// It provides helpers to:
//   - register unauthenticated HTTP transport routes that receive CloudEvents and dispatch to a ReceiverFunc, and
//   - create CloudEvents HTTP clients that send events using an HTTP RoundTripper and configured encoding.
//
// Receiver panics are recovered at the CloudEvents callback boundary. The CloudEvents SDK otherwise recovers
// callback panics internally, before go-service's transport recovery middleware can write its safe response or
// emit its request log. A recovered panic becomes a safe HTTP 500/NACK response, and its diagnostic error is
// recorded for the standard HTTP logger without being exposed to the sender.
//
// Start with [NewReceiver] / [Receiver.Register] for receiving and [NewSender] for sending.
package events
