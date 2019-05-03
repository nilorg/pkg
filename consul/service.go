package consul

import (
	"context"
	"log"

	"github.com/nilorg/pkg/consul/health/grpc_health_v1"
	"google.golang.org/grpc"
)

var (
	// debug mode
	Debug = false
)

// HealthService 健康检查服务
type HealthService struct {
	ServiceName string
}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthService) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if Debug {
		log.Printf("%s Service Check...", req.GetService())
	}
	status := grpc_health_v1.HealthCheckResponse_SERVING
	if h.ServiceName != req.Service {
		status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}

// RegisterHealthServer 注册健康检查服务
func RegisterHealthServer(s *grpc.Server, serviceName string) {
	grpc_health_v1.RegisterHealthServer(s, &HealthService{ServiceName: serviceName})
}
