package tls

import "crypto/tls"

// Config is an alias of crypto/tls.Config.
//
// It configures TLS handshake behavior, certificates, protocol versions, and
// peer verification for clients and servers.
//
// The alias exists so repository packages can stay on the go-service import
// path without changing the behavior or shape of the standard library type.
type Config = tls.Config

// Certificate is an alias of crypto/tls.Certificate.
//
// It holds a parsed certificate chain and corresponding private key used during
// TLS handshakes.
//
// The alias preserves the standard library representation exactly.
type Certificate = tls.Certificate

// VersionTLS12 is an alias of crypto/tls.VersionTLS12.
//
// It identifies TLS 1.2 for use in runtime TLS config defaults and protocol
// minimum-version checks.
const VersionTLS12 = tls.VersionTLS12

// RequireAndVerifyClientCert is an alias of
// crypto/tls.RequireAndVerifyClientCert.
//
// It configures mutual TLS by requiring a client certificate and verifying it
// during the handshake.
const RequireAndVerifyClientCert = tls.RequireAndVerifyClientCert

// X509KeyPair parses a public/private key pair from PEM data.
//
// This is a thin wrapper around crypto/tls.X509KeyPair and preserves the
// standard library behavior exactly, including PEM parsing and validation
// errors.
func X509KeyPair(certPEMBlock, keyPEMBlock []byte) (Certificate, error) {
	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}
