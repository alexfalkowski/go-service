package method_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/grpc/method"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestPolicyOperationService(t *testing.T) {
	policy := method.NewPolicy()
	desc := &grpc.ServiceDesc{
		ServiceName: "greet.v1.GreeterService",
		Methods: []grpc.MethodDesc{
			{MethodName: "SayHello"},
		},
		Streams: []grpc.StreamDesc{
			{StreamName: "SayStreamHello"},
		},
	}

	policy.OperationService(desc)

	require.True(t, policy.IsOperation("/greet.v1.GreeterService/SayHello"))
	require.True(t, policy.IsOperation("/greet.v1.GreeterService/SayStreamHello"))
	require.False(t, policy.IsOperation("/greet.v1.GreeterService/Other"))
}

func TestPolicyAllowUnauthenticatedService(t *testing.T) {
	policy := method.NewPolicy()
	desc := &grpc.ServiceDesc{
		ServiceName: "events.v1.EventsService",
		Methods: []grpc.MethodDesc{
			{MethodName: "Receive"},
		},
		Streams: []grpc.StreamDesc{
			{StreamName: "Subscribe"},
		},
	}

	policy.AllowUnauthenticatedService(desc)

	require.True(t, policy.IsUnauthenticated("/events.v1.EventsService/Receive"))
	require.True(t, policy.IsUnauthenticated("/events.v1.EventsService/Subscribe"))
	require.False(t, policy.IsUnauthenticated("/events.v1.EventsService/Other"))
	require.False(t, policy.IsOperation("/events.v1.EventsService/Receive"))
}
