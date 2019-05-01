package consul

import (
	"fmt"
	"testing"
	"time"
)

// newRegisterInfo ...
func newRegisterInfo() *RegisterInfo {
	ip := LocalIP()
	fmt.Printf("服务端IP:%s\n ", ip)
	return &RegisterInfo{
		ServiceName:                    "unknown",
		IP:                             ip,
		Tags:                           []string{},
		Port:                           3000,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(10) * time.Second,
	}
}

func TestRegisterService(t *testing.T) {
	info := newRegisterInfo()
	client := NewClient("127.0.0.1:8500")
	err := client.Register(info)
	if err != nil {
		fmt.Printf("TestRegisterService :%v", err)
		t.Fatalf("TestRegisterService :%v", err)
		return
	}
}
