package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/nilorg/sdk/convert"

	"github.com/nilorg/pkg/consul"
	"google.golang.org/grpc"
)

const (
	serviceName      = "unknown"
	consulServerAddr = "127.0.0.1:8500"
)

// 用于测试consul的服务注册

func main() {
	port := randInt64(3000, 3500)
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register reflection service on gRPC server.
	s := grpc.NewServer()
	consul.RegisterHealthServer(s, serviceName)
	fmt.Printf("grpc addr %s\n", addr)
	rinfo := consul.NewDefaultRegisterInfo()
	rinfo.IP = consul.LocalIP()
	rinfo.Port = convert.ToInt(port)
	consulClient := consul.NewClient(consulServerAddr, serviceName)
	consulClient.Register(rinfo)
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

func randInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
