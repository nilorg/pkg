package zlog

import (
	"context"

	"github.com/nilorg/sqlxplus"
	"go.uber.org/zap"
)

// SqlxplusLogger ...
type SqlxplusLogger struct {
}

// Printf 打印
func (SqlxplusLogger) Printf(ctx context.Context, query string, args ...interface{}) {
	indexs := sqlxplus.StringIndex(query, '?')
	query = sqlxplus.StringIndexReplace(query, indexs, args)
	WithSugared(ctx).WithOptions(zap.AddCallerSkip(3)).Debugf("[sqlx] %s", query)
}

// Println 打印
func (SqlxplusLogger) Println(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[sqlx]",
	}
	nArgs = append(nArgs, args...)
	WithSugared(ctx).WithOptions(zap.AddCallerSkip(3)).Debug(nArgs...)
}

// Errorf 错误
func (SqlxplusLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	WithSugared(ctx).WithOptions(zap.AddCallerSkip(3)).Errorf("[sqlx-error] "+format, args...)
}

// Errorln 错误
func (SqlxplusLogger) Errorln(ctx context.Context, args ...interface{}) {
	nArgs := []interface{}{
		"[sqlx-error]",
	}
	nArgs = append(nArgs, args...)
	WithSugared(ctx).WithOptions(zap.AddCallerSkip(3)).Error(nArgs...)
}
