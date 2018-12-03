package grpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

// Server 服务端
type Server struct {
	address   string
	rpcServer *grpc.Server
}

// NewServer 创建服务端
func NewServer(address string, tls bool, certFile, keyFile string) *Server {
	var rpcServer *grpc.Server
	if tls {
		// TLS认证
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		// 实例化grpc Server, 并开启TLS认证
		rpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		rpcServer = grpc.NewServer()
	}

	return &Server{
		rpcServer: rpcServer,
		address:   address,
	}
}

// Start 启动
func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		grpclog.Fatalf("grpc failed to listen: %v", err)
	}
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.rpcServer)
	go func() {
		if err := s.rpcServer.Serve(lis); err != nil {
			grpclog.Fatalf("grpc failed to serve: %v", err)
		}
	}()
}

// GetSrv 获取rpc server
func (s *Server) GetSrv() *grpc.Server {
	return s.rpcServer
}

// Close 关闭
func (s *Server) Close() {
	s.rpcServer.Stop()
}