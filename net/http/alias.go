package http

import (
	"net/http"

	"github.com/urfave/negroni/v3"
)

// DefaultMaxHeaderBytes is an alias of [http.DefaultMaxHeaderBytes].
const DefaultMaxHeaderBytes = http.DefaultMaxHeaderBytes

// TimeFormat is the HTTP time format layout.
//
// It is an alias of [http.TimeFormat].
const TimeFormat = http.TimeFormat

// Client is an alias for [net/http.Client].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Client = http.Client

// MaxBytesError is an alias for [net/http.MaxBytesError].
//
// It is returned when MaxBytesReader or MaxBytesHandler observes an inbound request body exceeding
// the configured byte limit.
type MaxBytesError = http.MaxBytesError

// Handler is an alias for [net/http.Handler].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Handler = http.Handler

// ChainedHandler is an alias for negroni.Handler.
type ChainedHandler = negroni.Handler

// ChainedHandlers is an alias for negroni.Negroni.
type ChainedHandlers = negroni.Negroni

// HandlerFunc is an alias for [net/http.HandlerFunc].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type HandlerFunc = http.HandlerFunc

// Header is an alias for [net/http.Header].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Header = http.Header

// Request is an alias for [net/http.Request].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Request = http.Request

// Response is an alias for [net/http.Response].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Response = http.Response

// ServeMux is an alias for [net/http.ServeMux].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ServeMux = http.ServeMux

// Server is an alias for [net/http.Server].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Server = http.Server

// ResponseWriter is an alias for [net/http.ResponseWriter].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ResponseWriter = http.ResponseWriter

// Flusher is an alias for [net/http.Flusher].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Flusher = http.Flusher

// ResponseController is an alias for [net/http.ResponseController].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ResponseController = http.ResponseController

// RoundTripper is an alias for [net/http.RoundTripper].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type RoundTripper = http.RoundTripper

// DefaultTransport is an alias for [http.DefaultTransport].
var DefaultTransport = http.DefaultTransport

// ErrUseLastResponse is an alias for [http.ErrUseLastResponse].
var ErrUseLastResponse = http.ErrUseLastResponse

// ErrServerClosed is an alias for [http.ErrServerClosed].
var ErrServerClosed = http.ErrServerClosed

// ErrAbortHandler is an alias for [http.ErrAbortHandler].
var ErrAbortHandler = http.ErrAbortHandler

// NoBody is an alias for [http.NoBody].
var NoBody = http.NoBody
