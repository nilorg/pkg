package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Client grpc客户端
type Client struct {
	conn          *grpc.ClientConn // 连接
	serverAddress string
	tls           bool
}

// GetConn 获取客户端连接
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// NewClient 创建grpc客户端
func NewClient(serverAddress string) *Client {
	return newClient(serverAddress, nil, nil)
}

// NewClientTLSFromFile 创建grpc客户端TLSFromFile
func NewClientTLSFromFile(serverAddress string, certFile, serverNameOverride string) *Client {
	// TLS连接
	creds, err := credentials.NewClientTLSFromFile(certFile, serverNameOverride)
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	return NewClientTLS(serverAddress, creds)
}

// NewClientTLS 创建grpc客户端
func NewClientTLS(serverAddress string, creds credentials.TransportCredentials) *Client {
	return newClient(serverAddress, creds, nil)
}

// CustomCredential 自定义凭证
type CustomCredential struct {
	AppKey, AppSecret string
	Security          bool
}

// NewCustomCredential 创建自定义凭证
func NewCustomCredential(appKey, appSecret string, tls bool) *CustomCredential {
	return &CustomCredential{
		AppKey:    appKey,
		AppSecret: appSecret,
		Security:  tls,
	}
}

// GetRequestMetadata Get请求元数据
func (c CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_key":    c.AppKey,
		"app_secret": c.AppSecret,
	}, nil
}

// RequireTransportSecurity 是否安全传输
func (c CustomCredential) RequireTransportSecurity() bool {
	return c.Security
}

// GetCustomAuthenticationParameter 获取自定义参数
type GetCustomAuthenticationParameter func() (appID, appKey string)

// NewClientCustomAuthentication 创建grpc客户端自定义服务验证
func NewClientCustomAuthentication(serverAddress string, credential credentials.PerRPCCredentials) *Client {
	return newClient(serverAddress, nil, credential)
}

// NewClientTLSCustomAuthentication 创建grpc客户端TLS自定义服务验证
func NewClientTLSCustomAuthentication(serverAddress string, creds credentials.TransportCredentials, credential credentials.PerRPCCredentials) *Client {
	return newClient(serverAddress, creds, credential)
}

func newClient(serverAddress string, creds credentials.TransportCredentials, credential credentials.PerRPCCredentials) *Client {
	var opts []grpc.DialOption
	if creds == nil && credential == nil {
		opts = append(opts, grpc.WithInsecure())
	} else {
		if creds != nil {
			opts = append(opts, grpc.WithInsecure())
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
		if credential != nil {
			// 使用自定义认证
			opts = append(opts, grpc.WithPerRPCCredentials(credential))
		}
	}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		grpclog.Fatalln(err)
	}
	return &Client{
		conn:          conn,
		serverAddress: serverAddress,
		tls:           creds != nil,
	}
}

// Close 关闭
func (c *Client) Close() {
	c.conn.Close()
}
