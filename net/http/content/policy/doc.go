// Package policy classifies HTTP content codecs for untrusted request decoding.
//
// Unary and streaming handlers use this package after resolving a registered codec. CanDecode rejects
// codecs without ratio and nesting bounds while leaving codec registration and response encoding unchanged.
package policy
