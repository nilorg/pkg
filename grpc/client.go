package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Client grpc客户端
type Client struct {
	conn               *grpc.ClientConn // 连接
	serverAddress      string
	tls                bool
	certFile           string
	serverNameOverride string
}

// GetConn 获取客户端连接
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// NewClient 创建grpc客户端
func NewClient(serverAddress string) *Client {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	return &Client{
		conn:               conn,
		serverAddress:      serverAddress,
		tls:                false,
		certFile:           "",
		serverNameOverride: "",
	}
}

// NewClientTLS 创建grpc客户端TLS
func NewClientTLS(serverAddress string, tls bool, certFile, serverNameOverride string) *Client {
	var conn *grpc.ClientConn
	var err error
	// TLS连接
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(certFile, serverNameOverride)
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err = grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))

	if err != nil {
		grpclog.Fatalf("did not connect: %v", err)
	}
	return &Client{
		conn:               conn,
		serverAddress:      serverAddress,
		tls:                tls,
		certFile:           certFile,
		serverNameOverride: serverNameOverride,
	}
}

// CustomCredential 自定义凭证
type CustomCredential struct {
	AppID  string
	AppKey string
}

// NewCustomCredential 创建自定义凭证
func NewCustomCredential(appID, appKey string) *CustomCredential {
	return &CustomCredential{
		AppID:  appID,
		AppKey: appKey,
	}
}

// GetRequestMetadata Get请求元数据
func (c CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  c.AppID,
		"app_key": c.AppKey,
	}, nil
}

// RequireTransportSecurity 是否安全传输
func (c CustomCredential) RequireTransportSecurity() bool {
	return false
}

// GetCustomAuthenticationParameter 获取自定义参数
type GetCustomAuthenticationParameter func() (appID, appKey string)

// NewClientCustomAuthentication 创建grpc客户端自定义服务验证
func NewClientCustomAuthentication(serverAddress string, credential credentials.PerRPCCredentials) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// 使用自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(credential))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		grpclog.Fatalln(err)
	}
	return &Client{
		conn:               conn,
		serverAddress:      serverAddress,
		tls:                false,
		certFile:           "",
		serverNameOverride: "",
	}
}

// Close 关闭
func (c *Client) Close() {
	c.conn.Close()
}
