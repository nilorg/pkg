package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nilorg/sdk/mq"

	"github.com/nats-io/nats.go"
)

// Nats ...
type Nats struct {
	conn *nats.Conn
}

//Publish ...
func (n *Nats) Publish(ctx context.Context, subj string, msg interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(msg)
	if err != nil {
		return
	}
	return n.conn.Publish(subj, data)
}

//Subscribe ...
func (n *Nats) Subscribe(topic string, h mq.SubscribeHandler, queue ...string) (err error) {
	if len(queue) > 0 {
		_, err = n.conn.QueueSubscribe(topic, queue[0], func(msg *nats.Msg) {
			h(context.Background(), msg.Data)
		})
	} else {
		_, err = n.conn.Subscribe(topic, func(msg *nats.Msg) {
			h(context.Background(), msg.Data)
		})
	}
	return
}

// NewNats ..
func NewNats(url string) (*Nats, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Nats{
		conn: conn,
	}, nil
}

// NewNatsOptions ..
func NewNatsOptions(url string) (*Nats, error) {
	opts := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   5,
		ReconnectWait:  5 * time.Second,
		Timeout:        3 * time.Second,
		Url:            url,
	}
	conn, err := opts.Connect()
	if err != nil {
		return nil, err
	}
	return &Nats{
		conn: conn,
	}, nil
}
