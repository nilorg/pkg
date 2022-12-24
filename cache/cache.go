package cache

import (
	"context"
	"time"
)

// Cacher 缓存
type Cacher interface {
	Geter
	Seter
	Remover
	Keyer
}

// Geter ...
type Geter interface {
	Get(ctx context.Context, key string, v interface{}) (err error)
	GetString(ctx context.Context, key string) (value string, err error)
	GetBytes(ctx context.Context, key string) (bytes []byte, err error)
	GetHash(ctx context.Context, key, field string) (value string)
}

// Seter ...
type Seter interface {
	Set(ctx context.Context, key string, v interface{}, expiration ...time.Duration) (err error)
	SetString(ctx context.Context, key string, value string, expiration ...time.Duration) (err error)
	SetBytes(ctx context.Context, key string, bytes []byte, expiration ...time.Duration) (err error)
	SetHash(ctx context.Context, key string, v map[string]string, expiration ...time.Duration) (err error)
}

// Remover ...
type Remover interface {
	Remove(ctx context.Context, keys ...string) (err error)
	RemoveMatch(ctx context.Context, match string) (err error)
}

type Keyer interface {
	// Zrem
	PushRangeKey(ctx context.Context, setKey string, keys ...string) (err error)
	DelRangeKey(ctx context.Context, setKey string) (err error)
}

type QueryFunc func(v interface{}) error

type Cache interface {
	Get(ctx context.Context, key string, v interface{}) (err error)
	Task(ctx context.Context, key string, v interface{}, query QueryFunc, expiry ...time.Duration) (err error)
	Set(ctx context.Context, key string, v interface{}, expiry ...time.Duration) (err error)
	Del(ctx context.Context, keys ...string) (err error)
	AddRangeKey(ctx context.Context, setKey string, keys ...string) (err error)
	DelRangeKey(ctx context.Context, setKey string) (err error)
	IsNotFound(err error) bool
}
