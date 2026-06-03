package log

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HTTP 与 Gin Context 中请求 ID 的键名
const (
	HeaderRequestID = "X-Request-ID"
	CtxRequestID    = "request_id"
)

// WithGin 在日志字段中附加上下文：request_id、http_method、path
func WithGin(c *gin.Context) *zap.Logger {
	rid := ""
	if v, ok := c.Get(CtxRequestID); ok {
		if s, ok := v.(string); ok {
			rid = s
		}
	}
	return L().With(
		zap.String("request_id", rid),
		zap.String("http_method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)
}
