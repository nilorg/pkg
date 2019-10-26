package mq

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

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
	err = client.Subscribe("test", func(ctx context.Context, msg interface{}) {
		log.Println(string(msg.([]byte)))
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
