package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/log"
	"go.uber.org/zap"
)

// RecoveryZap 捕获 panic，用 zap 打日志（含 request_id）并返回 500
func RecoveryZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.FromContext(c.Request.Context()).Error("panic",
					zap.Any("recover", rec),
					zap.String("error_type", "panic"),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
