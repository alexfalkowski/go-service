// Package events provides CloudEvents HTTP sender/receiver wiring for go-service.
//
// This package integrates go-service CloudEvents HTTP helpers with transport wiring.
// It provides helpers to:
//   - register HTTP handlers that receive CloudEvents and dispatch to a ReceiverFunc, and
//   - create CloudEvents HTTP clients that send events using an HTTP RoundTripper.
//
// Start with [NewReceiver] / [Receiver.Register] for receiving and [NewSender] for sending.
package events
