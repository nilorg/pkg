package cache

import (
	"context"
)

type CacheKey struct{}

// NewCacheContext ...
func NewCacheContext(ctx context.Context, cache Cacher) context.Context {
	return context.WithValue(ctx, CacheKey{}, cache)
}

// FromCacheContext ...
func FromCacheContext(ctx context.Context) (cache Cacher, ok bool) {
	cache, ok = ctx.Value(CacheKey{}).(Cacher)
	return
}

type SkipCacheKey struct{}

// NewSkipCacheContext 创建跳过缓存到上下文
func NewSkipCacheContext(ctx context.Context, skip ...bool) context.Context {
	s := true
	if len(skip) > 0 {
		s = skip[0]
	}
	return context.WithValue(ctx, SkipCacheKey{}, s)
}

// FromSkipCacheContext 从上下文中获取跳过缓存变量
func FromSkipCacheContext(ctx context.Context) (skip bool) {
	var ok bool
	skip, ok = ctx.Value(SkipCacheKey{}).(bool)
	if !ok {
		skip = false
	}
	return
}
