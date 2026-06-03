package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

var loggerKey ctxKey

// NewContext 将请求级 Logger 写入 context，供 service/dao 通过 FromContext 取用
func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext 从 context 取出请求级 Logger；未注入时回退到根 Logger
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok && l != nil {
		return l
	}
	return L()
}
