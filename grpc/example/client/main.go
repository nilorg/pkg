package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nilorg/pkg/grpc"
	pb "github.com/nilorg/pkg/grpc/example/helloworld"
)

func main() {
	client := grpc.NewClient("127.0.0.1:5000")
	greeterClient := pb.NewGreeterClient(client.GetConn())

	go func() {
		for {
			time.Sleep(time.Second)
			r, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: "xudeyi"})
			if err != nil {
				log.Printf("could not greet: %v", err)
				continue
			}
			log.Printf("Greeting: %s", r.Message)

		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	defer client.Close()
}
