package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nilorg/pkg/grpc"
	pb "github.com/nilorg/pkg/grpc/example/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	ser := grpc.NewServer(":5000")
	ser.Start()
	pb.RegisterGreeterServer(ser.GetSrv(), &server{})
	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	defer ser.Close()
}
