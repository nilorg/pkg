package consul

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/naming"
)

// Clienter ...
type Clienter interface {
	Register(rInfo *RegisterInfo) error
	Resolve(target string) (naming.Watcher, error)
}

// NewClient ...
func NewClient(serverAddress, serviceName string, tags ...string) Clienter {
	return &Client{
		ConsulServerAddress: serverAddress,
		ServiceName:         serviceName,
		Tags:                tags,
	}
}

// Client consul客户端
type Client struct {
	ConsulServerAddress string
	ServiceName         string
	Tags                []string
}

// RegisterInfo consul service register info.
type RegisterInfo struct {
	IP                             string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

// NewDefaultRegisterInfo 创建默认注册信息
func NewDefaultRegisterInfo() *RegisterInfo {
	return &RegisterInfo{
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

// Register ...
func (c *Client) Register(rInfo *RegisterInfo) error {
	config := api.DefaultConfig()
	config.Address = c.ConsulServerAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()
	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", c.ServiceName, rInfo.IP, rInfo.Port), // 服务节点的名称
		Name:    c.ServiceName,                                                // 服务名称
		Tags:    c.Tags,                                                       // tag，可以为空
		Port:    rInfo.Port,                                                   // 服务端口
		Address: rInfo.IP,                                                     // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       rInfo.Interval.String(),                                      // 健康检查间隔
			GRPC:                           fmt.Sprintf("%v:%v/%v", rInfo.IP, rInfo.Port, c.ServiceName), // grpc 支持，执行健康检查的地址，ServiceName 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: rInfo.DeregisterCriticalServiceAfter.String(),                // 注销时间，相当于过期时间
		},
	}
	return agent.ServiceRegister(reg)
}

// Resolve ...
func (c *Client) Resolve(target string) (naming.Watcher, error) {
	config := api.DefaultConfig()
	config.Address = c.ConsulServerAddress
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulWatcher{
		client:  client,
		service: c.ServiceName,
		addrs:   map[string]struct{}{},
		tags:    c.Tags,
	}, nil
}

type consulWatcher struct {
	client    *api.Client
	service   string
	addrs     map[string]struct{}
	lastIndex uint64
	tags      []string
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	for {
		services, metainfo, err := w.client.Health().ServiceMultipleTags(w.service, w.tags, true, &api.QueryOptions{
			WaitIndex: w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})
		if err != nil {
			logrus.Warnf("error retrieving instances from Consul: %v", err)
		}
		w.lastIndex = metainfo.LastIndex

		addrs := map[string]struct{}{}
		for _, service := range services {
			addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
		}

		var updates []*naming.Update
		for addr := range w.addrs {
			if _, ok := addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr})
			}
		}

		for addr := range addrs {
			if _, ok := w.addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
			}
		}

		if len(updates) != 0 {
			w.addrs = addrs
			return updates, nil
		}
	}
}

func (w *consulWatcher) Close() {
	// nothing to do
}
