package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/nilorg/sdk/mq"
)

// Stan ...
type Stan struct {
	conn    stan.Conn
	subOpts []stan.SubscriptionOption
}

//Publish ...
func (n *Stan) Publish(ctx context.Context, subj string, msg interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(msg)
	if err != nil {
		return
	}
	return n.conn.Publish(subj, data)
}

//Subscribe ...
func (n *Stan) Subscribe(topic string, h mq.SubscribeHandler, queue ...string) (err error) {
	if len(queue) > 0 {
		_, err = n.conn.QueueSubscribe(topic, queue[0], func(msg *stan.Msg) {
			if msgErr := h(context.Background(), msg.Data); msgErr == nil {
				msg.Ack()
			}
		}, n.subOpts...)
	} else {
		_, err = n.conn.Subscribe(topic, func(msg *stan.Msg) {
			if msgErr := h(context.Background(), msg.Data); msgErr == nil {
				msg.Ack()
			}
		}, n.subOpts...)
	}
	return
}

// NewStan ..
func NewStan(conn stan.Conn, subOpts ...stan.SubscriptionOption) (*Stan, error) {
	if len(subOpts) == 0 {
		subOpts = append(subOpts, stan.SetManualAckMode(), stan.AckWait(time.Second*10))
	}
	return &Stan{
		conn:    conn,
		subOpts: subOpts,
	}, nil
}
