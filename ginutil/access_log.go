package ginutil

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/log"
	"go.uber.org/zap"
)

// AccessLog 使用结构化日志记录请求与响应，每条含 request_id，便于按请求聚类
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}
		c.Next()
		cost := time.Since(start)
		log.FromContext(c.Request.Context()).Info("access",
			zap.Int("http_status", c.Writer.Status()),
			zap.String("http_client_ip", c.ClientIP()),
			zap.Int64("latency_ms", cost.Milliseconds()),
			zap.String("path", path),
		)
	}
}
