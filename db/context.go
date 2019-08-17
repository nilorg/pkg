package db

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

var (
	// ErrContextNotFoundGorm 上下文不存在Gorm错误
	ErrContextNotFoundGorm = errors.New("上下文中没有获取到Gorm")
)

type gormKey struct{}

// FromContext 从上下文中获取微信客户端
func FromContext(ctx context.Context) (*gorm.DB, error) {
	c, ok := ctx.Value(gormKey{}).(*gorm.DB)
	if !ok {
		return nil, ErrContextNotFoundGorm
	}
	return c, nil
}

// NewContext 创建微信客户端上下文
func NewContext(ctx context.Context, gdb *gorm.DB) context.Context {
	return context.WithValue(ctx, gormKey{}, gdb)
}
