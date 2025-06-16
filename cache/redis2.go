package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nilorg/sync"
)

// RedisCache2 redis缓存
type RedisCache2 struct {
	redisClient *redis.Client
	opts        Options
	redisSync   *sync.RedisSync
}

// NewRedisCache2 ...
func NewRedisCache2(redisClient *redis.Client, opts ...Option) *RedisCache2 {
	return &RedisCache2{
		redisClient: redisClient,
		opts:        newOptions(opts...),
		redisSync:   sync.NewRedisSync(redisClient),
	}
}

func (r *RedisCache2) getBytes(ctx context.Context, key string) (bytes []byte, err error) {
	return r.redisClient.Get(ctx, key).Bytes()
}

func (r *RedisCache2) setBytes(ctx context.Context, key string, bytes []byte, expiry ...time.Duration) (err error) {
	if len(expiry) > 0 {
		err = r.redisClient.Set(ctx, key, bytes, expiry[0]).Err()
	} else {
		err = r.redisClient.Set(ctx, key, bytes, r.opts.Expiry()).Err()
	}
	return
}

func (r *RedisCache2) Get(ctx context.Context, key string, v interface{}) (err error) {
	key = r.formatKey(key)
	var bytes []byte
	if bytes, err = r.getBytes(ctx, key); err != nil {
		if r.IsNotFound(err) {
			_ = r.setBytes(ctx, key, []byte("CACHE_NOT_FOUND"), r.opts.NotFoundExpiry)
			err = r.opts.ErrNotFound
		}
		return
	}
	if string(bytes) == "CACHE_NOT_FOUND" {
		err = r.opts.ErrNotFound
		return
	}
	err = r.opts.Serializer.Unmarshal(bytes, v)
	return
}

func (r *RedisCache2) Task(ctx context.Context, key string, v interface{}, query QueryFunc, expiry ...time.Duration) (err error) {
	if FromSkipCacheContext(ctx) {
		err = query(v)
		return
	}
	err = r.Get(ctx, key, v)
	if r.IsNotFound(err) {
		// 一个存在的key，在缓存过期的瞬间，同时有大量的请求过来，造成所有请求都去读dB，这些请求都会击穿到DB，造成瞬时DB请求量大、压力骤增。
		mtx := r.redisSync.NewMutex(
			fmt.Sprintf("%s:%s", key, "lock"),
			sync.KeyPrefix(r.opts.Prefix),
		)
		if err = mtx.Lock(); err != nil {
			return
		}
		defer mtx.Unlock() // 使用defer确保锁被释放

		// 再次检查缓存，避免重复查询
		if err = r.Get(ctx, key, v); err == nil {
			return nil // 缓存已存在，直接返回
		}
		if !r.IsNotFound(err) {
			return err // 其他错误，直接返回
		}

		// 执行查询
		if err = query(v); err != nil {
			return err
		}

		// 设置缓存
		if err = r.Set(ctx, key, v, expiry...); err != nil {
			return err
		}
	}
	return nil
}
func (r *RedisCache2) Set(ctx context.Context, key string, v interface{}, expiry ...time.Duration) (err error) {
	var bytes []byte
	bytes, err = r.opts.Serializer.Marshal(v)
	if err != nil {
		return
	}
	key = r.formatKey(key)
	err = r.setBytes(ctx, key, bytes, expiry...)
	return
}

func (r *RedisCache2) Del(ctx context.Context, keys ...string) (err error) {
	l := len(keys)
	for i := 0; i < l; i++ {
		keys[i] = r.formatKey(keys[i])
	}
	err = r.redisClient.Del(ctx, keys...).Err()
	return
}

func (r *RedisCache2) formatKey(key string) string {
	return r.opts.Prefix + key
}

func (r *RedisCache2) AddRangeKey(ctx context.Context, setKey string, keys ...string) (err error) {
	setKey = r.formatKey(setKey)
	kl := len(keys)
	if kl == 0 {
		return
	}
	members := make([]interface{}, kl)
	for i := 0; i < kl; i++ {
		members[i] = r.formatKey(keys[i])
	}
	err = r.redisClient.SAdd(ctx, setKey, members...).Err()
	return
}

func (r *RedisCache2) DelRangeKey(ctx context.Context, setKey string) (err error) {
	setKey = r.formatKey(setKey)
	var keys []string
	if keys, err = r.redisClient.SMembers(ctx, setKey).Result(); err != nil {
		return
	}
	if len(keys) == 0 {
		return
	}
	err = r.redisClient.Del(ctx, keys...).Err()
	return
}

func (r *RedisCache2) IsNotFound(err error) bool {
	return errors.Is(err, redis.Nil) || errors.Is(err, r.opts.ErrNotFound)
}
