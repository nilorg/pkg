package hooks

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// RedisLogrusHook ...
type RedisLogrusHook struct {
	Channel   string
	redis     *redis.Client
	formatter logrus.Formatter
	LogLevels []logrus.Level
}

// NewRedisLogrusHook ...
func NewRedisLogrusHook(redis *redis.Client, formatter logrus.Formatter, channel string) logrus.Hook {
	return &RedisLogrusHook{
		redis:     redis,
		formatter: formatter,
		Channel:   channel,
		LogLevels: logrus.AllLevels,
	}
}

// Fire ...
func (h *RedisLogrusHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	if e.Context == nil {
		e.Context = context.Background()
	}
	err = h.redis.Publish(e.Context, h.Channel, dataBytes).Err()
	return err
}

// Levels returns all logrus levels.
func (h *RedisLogrusHook) Levels() []logrus.Level {
	return h.LogLevels
}
