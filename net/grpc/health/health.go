package health

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Request is an alias of the standard gRPC health check request.
type Request = health.HealthCheckRequest

// Response is an alias of the standard gRPC health check response.
type Response = health.HealthCheckResponse

// Status is an alias of the standard gRPC health serving status enum.
type Status = health.HealthCheckResponse_ServingStatus

// Client is an alias of the standard gRPC health client interface.
type Client = health.HealthClient

// ListRequest is an alias of the standard gRPC health list request.
type ListRequest = health.HealthListRequest

// ListResponse is an alias of the standard gRPC health list response.
type ListResponse = health.HealthListResponse

// Server is an alias of the standard gRPC health server interface.
type Server = health.HealthServer

// WatchClient is an alias of the standard gRPC health watch client stream.
type WatchClient = health.Health_WatchClient

// WatchServer is an alias of the standard gRPC health watch server stream.
type WatchServer = health.Health_WatchServer

// UnimplementedServer is an alias of the standard forward-compatible health server embedding type.
type UnimplementedServer = health.UnimplementedHealthServer

// UnsafeServer is an alias of the standard opt-out interface for forward-compatible health servers.
type UnsafeServer = health.UnsafeHealthServer

// CheckFullMethodName is the full method name for health Check RPCs.
const CheckFullMethodName = health.Health_Check_FullMethodName

// ListFullMethodName is the full method name for health List RPCs.
const ListFullMethodName = health.Health_List_FullMethodName

// WatchFullMethodName is the full method name for health Watch RPCs.
const WatchFullMethodName = health.Health_Watch_FullMethodName

// Unknown aliases the UNKNOWN serving status.
const Unknown = health.HealthCheckResponse_UNKNOWN

// Serving aliases the SERVING status.
const Serving = health.HealthCheckResponse_SERVING

// NotServing aliases the NOT_SERVING status.
const NotServing = health.HealthCheckResponse_NOT_SERVING

// ServiceUnknown aliases the SERVICE_UNKNOWN status.
const ServiceUnknown = health.HealthCheckResponse_SERVICE_UNKNOWN

// FileDescriptor aliases the upstream health proto file descriptor.
var FileDescriptor = health.File_grpc_health_v1_health_proto

// ServiceDesc aliases the upstream health service descriptor.
var ServiceDesc = health.Health_ServiceDesc

// StatusName aliases the upstream serving status name map.
var StatusName = health.HealthCheckResponse_ServingStatus_name

// StatusValue aliases the upstream serving status value map.
var StatusValue = health.HealthCheckResponse_ServingStatus_value

// NewClient returns a standard gRPC health client.
func NewClient(cc grpc.ClientConnInterface) Client {
	return health.NewHealthClient(cc)
}

// RegisterServer registers srv as a standard gRPC health service.
func RegisterServer(s grpc.ServiceRegistrar, srv Server) {
	health.RegisterHealthServer(s, srv)
}
