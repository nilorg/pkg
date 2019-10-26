package mq

import (
	"context"

	"github.com/nilorg/sdk/mq"
)

// Publisher ...
type Publisher struct {
	topic  string
	client mq.Clienter
}

// Publish ...
func (p *Publisher) Publish(ctx context.Context, msg interface{}) error {
	return p.client.Publish(ctx, p.topic, msg)
}

// NewPublisher returns a new Publisher
func NewPublisher(topic string, client mq.Clienter) mq.Publisher {
	return &Publisher{
		topic:  topic,
		client: client,
	}
}
