package method

import "google.golang.org/grpc"

// NewPolicy constructs a method policy.
func NewPolicy() *Policy {
	return &Policy{
		operations:      map[string]struct{}{},
		unauthenticated: map[string]struct{}{},
	}
}

// Policy stores gRPC full-method behavior used by server middleware.
//
// Populate Policy during service registration before the server starts. See
// the package documentation for the middleware behavior of each classification.
type Policy struct {
	operations      map[string]struct{}
	unauthenticated map[string]struct{}
}

// Operation marks name as an operation method.
//
// Operation methods bypass transport token verification, access control, and
// logging. Unary operations also bypass metadata extraction and limiting,
// while operation streams retain metadata extraction and stream limiting.
func (p *Policy) Operation(name string) {
	p.operations[name] = struct{}{}
}

// AllowUnauthenticated marks name as not requiring transport token authentication.
//
// Unauthenticated methods also bypass access control, but retain the normal
// metadata extraction, logging, and limiting behavior.
func (p *Policy) AllowUnauthenticated(name string) {
	p.unauthenticated[name] = struct{}{}
}

// OperationService marks every unary and stream method in desc as an operation method.
//
// See [Policy.Operation] for the middleware consequences.
func (p *Policy) OperationService(desc *grpc.ServiceDesc) {
	for _, method := range desc.Methods {
		p.Operation(fullMethodName(desc.ServiceName, method.MethodName))
	}
	for _, stream := range desc.Streams {
		p.Operation(fullMethodName(desc.ServiceName, stream.StreamName))
	}
}

// AllowUnauthenticatedService marks every unary and stream method in desc as not requiring transport token authentication.
//
// See [Policy.AllowUnauthenticated] for the middleware consequences.
func (p *Policy) AllowUnauthenticatedService(desc *grpc.ServiceDesc) {
	for _, method := range desc.Methods {
		p.AllowUnauthenticated(fullMethodName(desc.ServiceName, method.MethodName))
	}
	for _, stream := range desc.Streams {
		p.AllowUnauthenticated(fullMethodName(desc.ServiceName, stream.StreamName))
	}
}

// IsOperation reports whether name is an operation method.
func (p *Policy) IsOperation(name string) bool {
	_, ok := p.operations[name]
	return ok
}

// IsUnauthenticated reports whether name does not require transport token authentication.
func (p *Policy) IsUnauthenticated(name string) bool {
	_, ok := p.unauthenticated[name]
	return ok
}

func fullMethodName(service, method string) string {
	return "/" + service + "/" + method
}
