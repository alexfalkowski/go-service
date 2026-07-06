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
type Policy struct {
	operations      map[string]struct{}
	unauthenticated map[string]struct{}
}

// Operation marks name as an operation method.
func (p *Policy) Operation(name string) {
	p.operations[name] = struct{}{}
}

// AllowUnauthenticated marks name as not requiring transport token authentication.
func (p *Policy) AllowUnauthenticated(name string) {
	p.unauthenticated[name] = struct{}{}
}

// OperationService marks every unary and stream method in desc as an operation method.
func (p *Policy) OperationService(desc *grpc.ServiceDesc) {
	for _, method := range desc.Methods {
		p.Operation(fullMethodName(desc.ServiceName, method.MethodName))
	}
	for _, stream := range desc.Streams {
		p.Operation(fullMethodName(desc.ServiceName, stream.StreamName))
	}
}

// AllowUnauthenticatedService marks every unary and stream method in desc as not requiring transport token authentication.
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
