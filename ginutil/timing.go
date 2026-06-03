package ginutil

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/log"
	"go.uber.org/zap"
)

const slowRequestThreshold = time.Second

// TimingStats 统计接口耗时；超过 1s 时使用 warn 级别
func TimingStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		l := log.FromContext(c.Request.Context())
		fields := []zap.Field{zap.Duration("用时", cost)}
		if cost > slowRequestThreshold {
			l.Warn("slow request", fields...)
		} else {
			l.Info("request timing", fields...)
		}
	}
}
