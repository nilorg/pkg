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
	userID       string
	redisClient  *redis.Client
	maxErrCount  int
	lockDuration time.Duration
}

// Option 配置选项函数类型
type Option func(*Login)

// WithMaxErrCount 设置最大错误次数
func WithMaxErrCount(count int) Option {
	return func(l *Login) {
		if count > 0 {
			l.maxErrCount = count
		}
	}
}

// WithLockDuration 设置锁定时长
func WithLockDuration(duration time.Duration) Option {
	return func(l *Login) {
		if duration > 0 {
			l.lockDuration = duration
		}
	}
}

func New(userID string, redisClient *redis.Client, opts ...Option) *Login {
	l := &Login{
		userID:       userID,
		redisClient:  redisClient,
		maxErrCount:  5,
		lockDuration: 24 * time.Hour,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
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
	// 先检查是否已锁定
	locked, err = a.IsLocked(ctx)
	if err != nil || locked {
		return
	}

	var count int64
	count, err = a.redisClient.Incr(ctx, a.countKey()).Result()
	if err != nil {
		return
	}
	if count == 1 {
		err = a.redisClient.Expire(ctx, a.countKey(), a.lockDuration).Err()
		if err != nil {
			return
		}
	}
	remainingCount = a.maxErrCount - int(count)
	if remainingCount <= 0 {
		remainingCount = 0
		err = a.redisClient.Set(ctx, a.lockKey(), lock, a.lockDuration).Err()
		if err != nil {
			return
		}
		locked = true
	}
	return
}

// Reset 重置错误次数
func (a *Login) Reset(ctx context.Context) (err error) {
	err = a.redisClient.Del(ctx, a.lockKey(), a.countKey()).Err()
	return
}

// GetMaxErrCount 获取最大错误次数
func (a *Login) GetMaxErrCount() int {
	return a.maxErrCount
}

// GetLockDuration 获取锁定时长
func (a *Login) GetLockDuration() time.Duration {
	return a.lockDuration
}
