package grpc_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/stretchr/testify/require"
)

func TestRegistrarRegisterService(t *testing.T) {
	serviceRegistrar := &stubServiceRegistrar{}
	policy := grpc.NewMethodPolicy()
	registrar := grpc.NewRegistrar(serviceRegistrar, policy)
	desc := &grpc.ServiceDesc{ServiceName: "greet.v1.GreeterService"}
	impl := new(struct{})

	registrar.RegisterService(desc, impl)

	require.Same(t, desc, serviceRegistrar.desc)
	require.Same(t, impl, serviceRegistrar.impl)
	require.False(t, policy.IsOperation("/greet.v1.GreeterService/SayHello"))
}

func TestRegistrarRegisterOperationService(t *testing.T) {
	serviceRegistrar := &stubServiceRegistrar{}
	policy := grpc.NewMethodPolicy()
	registrar := grpc.NewRegistrar(serviceRegistrar, policy)
	desc := &grpc.ServiceDesc{
		ServiceName: "grpc.health.v1.Health",
		Methods: []grpc.MethodDesc{
			{MethodName: "Check"},
			{MethodName: "List"},
		},
		Streams: []grpc.StreamDesc{
			{StreamName: "Watch"},
		},
	}
	impl := new(struct{})

	registrar.RegisterOperationService(desc, impl)

	require.Same(t, desc, serviceRegistrar.desc)
	require.Same(t, impl, serviceRegistrar.impl)
	require.True(t, policy.IsOperation("/grpc.health.v1.Health/Check"))
	require.True(t, policy.IsOperation("/grpc.health.v1.Health/List"))
	require.True(t, policy.IsOperation("/grpc.health.v1.Health/Watch"))
}

func TestRegistrarRegisterUnauthenticatedService(t *testing.T) {
	serviceRegistrar := &stubServiceRegistrar{}
	policy := grpc.NewMethodPolicy()
	registrar := grpc.NewRegistrar(serviceRegistrar, policy)
	desc := &grpc.ServiceDesc{
		ServiceName: "events.v1.EventsService",
		Methods: []grpc.MethodDesc{
			{MethodName: "Receive"},
		},
		Streams: []grpc.StreamDesc{
			{StreamName: "Subscribe"},
		},
	}
	impl := new(struct{})

	registrar.RegisterUnauthenticatedService(desc, impl)

	require.Same(t, desc, serviceRegistrar.desc)
	require.Same(t, impl, serviceRegistrar.impl)
	require.True(t, policy.IsUnauthenticated("/events.v1.EventsService/Receive"))
	require.True(t, policy.IsUnauthenticated("/events.v1.EventsService/Subscribe"))
	require.False(t, policy.IsOperation("/events.v1.EventsService/Receive"))
}

type stubServiceRegistrar struct {
	desc *grpc.ServiceDesc
	impl any
}

func (r *stubServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, impl any) {
	r.desc = desc
	r.impl = impl
}
