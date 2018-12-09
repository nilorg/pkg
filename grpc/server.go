package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var (
	// defaultServer 默认server
	defaultServer *Server
)

// Server 服务端
type Server struct {
	address string
	tls       bool
	certFile  string
	keyFile   string
	rpcServer *grpc.Server
}

// NewServer 创建服务端
func NewServer(address string) *Server {
	var rpcServer *grpc.Server
	rpcServer = grpc.NewServer()
	return &Server{
		rpcServer: rpcServer,
		address:   address,
		tls:       false,
		certFile:  "",
		keyFile:   "",
	}
}

// NewServerTLS 创建服务端TLS
func NewServerTLS(address string, tls bool, certFile, keyFile string) *Server {

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	// 实例化grpc Server, 并开启TLS认证
	rpcServer := grpc.NewServer(grpc.Creds(creds))

	return &Server{
		rpcServer: rpcServer,
		address:   address,
		tls:       tls,
		certFile:  certFile,
		keyFile:   keyFile,
	}
}

// ValidationFunc 验证方法
type ValidationFunc func(appID, appKey string) bool

// NewServerCustomAuthentication 创建服务端自定义服务验证
func NewServerCustomAuthentication(address string, validation ValidationFunc) *Server {
	rpcServer := grpc.NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "无令牌认证信息")
		}
		var (
			appID  string
			appKey string
		)

		if v, ok := md["app_id"]; ok {
			appID = v[0]
		}
		if v, ok := md["app_key"]; ok {
			appKey = v[0]
		}
		if !validation(appID, appKey) {
			return nil, grpc.Errorf(codes.Unauthenticated, "无效的认证信息")
		}
		return handler(ctx, req)
	}))
	return &Server{
		rpcServer: rpcServer,
		address:   address,
		tls:       false,
		certFile:  "",
		keyFile:   "",
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

// Start 启动Grpc
func Start(address string) {
	defaultServer = NewServer(address)
	defaultServer.Start()
}

// GetSrv 获取rpc server
func GetSrv() *grpc.Server {
	return defaultServer.GetSrv()
}

// Close 关闭Grpc
func Close() {
	defaultServer.Close()
}
