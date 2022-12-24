package cache

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var testCache Cache

func TestMain(m *testing.M) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	testCache = NewRedisCache2(redisClient)
	m.Run()
}

func TestRedis2_GetCacheNotFound(t *testing.T) {
	var obj interface{}
	err := testCache.Get(context.Background(), "test", &obj)
	assert.Equal(t, true, testCache.IsNotFound(err))
	assert.NotEqual(t, err, ErrDefaultNotFound, obj)
}

func TestRedis2_SetCache(t *testing.T) {
	err := testCache.Set(context.Background(), "test", "test")
	assert.NoError(t, err)
}

func TestRedis2_GetCache(t *testing.T) {
	var obj interface{}
	err := testCache.Get(context.Background(), "test", &obj)
	assert.NotEqual(t, err, ErrDefaultNotFound, obj)
}

func TestRedis2_DelCache(t *testing.T) {
	err := testCache.Del(context.Background(), "test")
	assert.NoError(t, err)
}
