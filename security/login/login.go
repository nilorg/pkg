package login

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// 1. 从 Redis 数据库中获取与该用户名相关联的错误登录次数。

// 2. 检查错误登录次数是否达到了锁定阈值。如果达到了阈值，将账户状态设置为锁定，并在 Redis 中设置过期时间。

// 3. 当用户尝试登录时，检查账户状态是否为锁定状态。

const (
	lock = "lock"
)

type Login struct {
	userID      string
	redisClient *redis.Client
	maxErrCount int
}

func New(userID string, redisClient *redis.Client) *Login {
	return &Login{
		userID:      userID,
		redisClient: redisClient,
		maxErrCount: 5,
	}
}

func (a *Login) lockKey() string {
	return fmt.Sprintf("security:login:user:%s:lock", a.userID)
}

func (a *Login) countKey() string {
	return fmt.Sprintf("security:login:user:%s:errcount", a.userID)
}

// IsLocked 判断是否锁定
func (a *Login) IsLocked(ctx context.Context) (locked bool, err error) {
	result := a.redisClient.Get(ctx, a.lockKey())
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			err = nil
			return
		}
		err = result.Err()
		return
	}
	if result.Val() == lock {
		locked = true
	}
	return
}

// TryLock 尝试锁定
func (a *Login) TryLock(ctx context.Context) (locked bool, remainingCount int, err error) {
	err = a.redisClient.Incr(ctx, a.countKey()).Err()
	if err != nil {
		return
	}
	var count int
	count, err = a.redisClient.Get(ctx, a.countKey()).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return
	}
	remainingCount = a.maxErrCount - count
	if remainingCount <= 0 {
		err = a.redisClient.Set(ctx, a.lockKey(), lock, 24*time.Hour).Err()
		if err != nil {
			return
		}
		locked = true
	}
	return
}

// Reset 重置错误次数
func (a *Login) Reset(ctx context.Context) (err error) {
	err = a.redisClient.Del(ctx, a.lockKey()).Err()
	if err != nil {
		return
	}
	err = a.redisClient.Del(ctx, a.countKey()).Err()
	if err != nil {
		return
	}
	return
}

// GetMaxErrCount 获取最大错误次数
func (a *Login) GetMaxErrCount() int {
	return a.maxErrCount
}
