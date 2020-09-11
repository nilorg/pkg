package mq

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"

	"github.com/nilorg/sdk/mq"
)

func TestNats(t *testing.T) {
	var client mq.Clienter
	var err error
	client, err = NewNats(nats.DefaultURL)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.Subscribe("test", func(ctx context.Context, data []byte) error {
		log.Println(string(data))
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	for {
		err = client.Publish(context.Background(), "test", "message ... "+time.Now().Format(time.RFC3339))
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(time.Second * 5)
	}

}

func TestStan(t *testing.T) {
	var (
		shopConn      stan.Conn
		shopClient    mq.Clienter
		crontabConn   stan.Conn
		crontabClient mq.Clienter
		err           error
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
	err = shopClient.Subscribe("test", func(ctx context.Context, data []byte) error {
		log.Println(string(data))
		return nil
	})
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
	for {
		err = crontabClient.Publish(context.Background(), "test", "message ... "+time.Now().Format(time.RFC3339))
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(time.Second * 5)
	}
}
