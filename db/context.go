package db

import (
	"context"
	"errors"

	gormV1 "github.com/jinzhu/gorm"
	gormV2 "gorm.io/gorm"
)

var (
	// ErrContextNotFoundGorm 上下文不存在Gorm错误
	ErrContextNotFoundGorm = errors.New("上下文中没有获取到Gorm")
)

type gormKey struct{}

// FromContext 从上下文中获取Gorm
func FromContext(ctx context.Context) (*gormV1.DB, error) {
	c, ok := ctx.Value(gormKey{}).(*gormV1.DB)
	if !ok {
		return nil, ErrContextNotFoundGorm
	}
	return c, nil
}

// NewContext 创建Gorm上下文
func NewContext(ctx context.Context, gdb *gormV1.DB) context.Context {
	return context.WithValue(ctx, gormKey{}, gdb)
}

type gormv2Key struct{}

// FromGormV2Context 从上下文中获取Gorm
func FromGormV2Context(ctx context.Context) (*gormV2.DB, error) {
	c, ok := ctx.Value(gormv2Key{}).(*gormV2.DB)
	if !ok {
		return nil, ErrContextNotFoundGorm
	}
	return c, nil
}

// NewGormV2Context 创建Gorm上下文
func NewGormV2Context(ctx context.Context, gdb *gormV2.DB) context.Context {
	return context.WithValue(ctx, gormv2Key{}, gdb)
}
