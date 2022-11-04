package zlog

import (
	"context"

	"github.com/nilorg/sqlxplus"
)

// SqlxplusLogger ...
type SqlxplusLogger struct {
}

// Printf 打印
func (SqlxplusLogger) Printf(ctx context.Context, query string, args ...interface{}) {
	indexs := sqlxplus.StringIndex(query, '?')
	query = sqlxplus.StringIndexReplace(query, indexs, args)
	WithSugared(ctx).Debugf("[sqlx] %s", query)
}

// Println 打印
func (SqlxplusLogger) Println(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[sqlx]",
	}
	nArgs = append(nArgs, args...)
	WithSugared(ctx).Debug(nArgs...)
}

// Errorf 错误
func (SqlxplusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	WithSugared(ctx).Errorf("[sqlx-error] "+format, args...)
}

// Errorln 错误
func (SqlxplusLogger) Errorln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[sqlx-error]",
	}
	nArgs = append(nArgs, args...)
	WithSugared(ctx).Error(nArgs...)
}
