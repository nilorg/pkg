package register

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

// ServiceInfo consul service register info.
type ServiceInfo struct {
	Name                           string
	Tags                           []string
	IP                             string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

// NewServiceInfo new consul service register info.
func NewServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

// Register consul 服务注册
func Register(consulServerAddress string, si *ServiceInfo) error {
	config := api.DefaultConfig()
	config.Address = consulServerAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", si.Name, si.Tags, si.Port), // 服务节点的名称
		Name:    si.Name,                                            // 服务名称
		Tags:    si.Tags,                                            // tag，可以为空
		Port:    si.Port,                                            // 服务端口
		Address: si.IP,                                              // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       si.Interval.String(),                             // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", si.IP, si.Port, si.Name), // grpc 支持，执行健康检查的地址，ServiceName 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: si.DeregisterCriticalServiceAfter.String(),       // 注销时间，相当于过期时间
		},
	}
	return agent.ServiceRegister(reg)
}
