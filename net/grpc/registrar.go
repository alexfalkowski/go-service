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

// RegisterOperationService marks desc as operations and registers it with the wrapped registrar.
func (r *Registrar) RegisterOperationService(desc *ServiceDesc, impl any) {
	r.policy.OperationService(desc)
	r.RegisterService(desc, impl)
}

// RegisterUnauthenticatedService marks desc as unauthenticated and registers it with the wrapped registrar.
func (r *Registrar) RegisterUnauthenticatedService(desc *ServiceDesc, impl any) {
	r.policy.AllowUnauthenticatedService(desc)
	r.RegisterService(desc, impl)
}
