package main

import (
	"fmt"

	"github.com/nilorg/pkg/consul"
)

// 用于测试consul的服务发现

func main() {

	client := consul.NewClient("127.0.0.1:8500")
	w, err := client.Resolve("unknown")
	if err != nil {
		fmt.Printf("client.Resolve:%v", err)
		return
	}
	ups, err := w.Next()
	if err != nil {
		fmt.Printf("w.Next():%v", err)
		return
	}
	for _, v := range ups {
		fmt.Printf("%+v\n", v)
	}
}
