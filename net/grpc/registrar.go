package grpc

// NewRegistrar constructs a gRPC registrar that can attach method policy while registering services.
func NewRegistrar(registrar ServiceRegistrar, policy *MethodPolicy) *Registrar {
	return &Registrar{registrar: registrar, policy: policy}
}

// Registrar wraps a gRPC service registrar with repository method policy.
type Registrar struct {
	registrar ServiceRegistrar
	policy    *MethodPolicy
}

// RegisterService registers desc with the wrapped registrar.
func (r *Registrar) RegisterService(desc *ServiceDesc, impl any) {
	r.registrar.RegisterService(desc, impl)
}

// RegisterOperationService marks every method in desc as an operation and registers the service.
//
// Operation methods bypass transport token verification, access control, and
// logging. Unary operations also bypass metadata extraction and limiting,
// while operation streams retain metadata extraction and stream limiting. Use
// this for operational services such as the standard gRPC health service, not
// for ordinary application RPCs.
func (r *Registrar) RegisterOperationService(desc *ServiceDesc, impl any) {
	r.policy.OperationService(desc)
	r.RegisterService(desc, impl)
}

// RegisterUnauthenticatedService marks every method in desc as unauthenticated and registers the service.
//
// These methods bypass token verification and access control while retaining
// normal metadata extraction, logging, and limiting.
func (r *Registrar) RegisterUnauthenticatedService(desc *ServiceDesc, impl any) {
	r.policy.AllowUnauthenticatedService(desc)
	r.RegisterService(desc, impl)
}
