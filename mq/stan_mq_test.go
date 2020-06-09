package mq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/nilorg/sdk/mq"
)

var (
	shopConn      stan.Conn
	shopClient    mq.Clienter
	crontabConn   stan.Conn
	crontabClient mq.Clienter
)

func TestStanShop(t *testing.T) {
	var (
		// shopConn      stan.Conn
		// shopClient    mq.Clienter
		// crontabConn   stan.Conn
		// crontabClient mq.Clienter
		err error
	)
	shopConn, err = stan.Connect("wohuitao", "shop2", stan.NatsURL("nats://nats-streaming.mq:4222"))
	if err != nil {
		t.Error(err)
		return
	}
	shopClient, err = NewStan(shopConn)
	if err != nil {
		t.Error(err)
		return
	}
	crontabConn, err = stan.Connect("wohuitao", "crontab", stan.NatsURL("nats://nats-streaming.mq:4222"))
	if err != nil {
		t.Error(err)
		return
	}
	crontabClient, err = NewStan(crontabConn)
	if err != nil {
		t.Error(err)
		return
	}

	// subscriber := mq.NewSubscriber(shopClient)
	// err = subscriber.Register("wohuitao.job.shop.topic.queue", func(ctx context.Context, data []byte) error {
	// 	fmt.Println("mq sub data: ", string(data))
	// 	return nil
	// }, "wohuitao.job.shop.topic.queue")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// go func() {
	// 	NewJob().RegisterSubscriber()
	// }()

	for {
		fmt.Println("1111111")
		publisher := mq.NewPublisher("wohuitao.job.shop.topic.queue", crontabClient)
		err := publisher.Publish(context.Background(), "message ... "+time.Now().Format(time.RFC3339))
		if err != nil {
			fmt.Println(err)
			return
		}
		// err = publisher.Publish(context.Background(), "wohuitao.job.shop.topic.queue", "message ... "+time.Now().Format(time.RFC3339))
		// if err != nil {
		// 	t.Error(err)
		// 	return
		// }
		time.Sleep(time.Second * 30)
	}
}

// Job ...
type Job struct {
}

// NewJob ...
func NewJob() *Job {
	return &Job{}
}

// RegisterExecuteHandler ...
func (j *Job) RegisterExecuteHandler(ctx context.Context, data []byte) (err error) {
	fmt.Println("mq sub data: ", string(data))
	return nil
}

// RegisterSubscriber 注册订阅
func (j *Job) RegisterSubscriber() {
	sub := mq.NewSubscriber(shopClient)
	err := sub.Register("wohuitao.job.shop.topic.queue", j.RegisterExecuteHandler, "wohuitao.job.shop.topic.queue")
	if err != nil {
		fmt.Printf("sub.Register Error: %s\n", err)
	}
}
