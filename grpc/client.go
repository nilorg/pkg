package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Client grpc客户端
type Client struct {
	conn *grpc.ClientConn // 连接
}

// GetConn 获取客户端连接
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// NewClient 创建grpc客户端
func NewClient(serverAddress string, tls bool, certFile, serverNameOverride string) *Client {
	var conn *grpc.ClientConn
	var err error
	if tls {
		// TLS连接
		var creds credentials.TransportCredentials
		creds, err = credentials.NewClientTLSFromFile(certFile, serverNameOverride)
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		conn, err = grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial(serverAddress, grpc.WithInsecure())
	}
	if err != nil {
		grpclog.Fatalf("did not connect: %v", err)
	}
	return &Client{
		conn: conn,
	}
}

// Close 关闭
func (c *Client) Close() {
	c.conn.Close()
}
