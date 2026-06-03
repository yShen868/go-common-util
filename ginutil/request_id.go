package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yShen868/go-common-util/log"
)

const maxHeaderRequestIDLen = 64

// RequestID 为每个请求生成/透传 X-Request-ID，写入 Gin 上下文，便于全链路日志关联
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(log.HeaderRequestID)
		if rid == "" || len(rid) > maxHeaderRequestIDLen {
			rid = uuid.NewString()
		}
		c.Set(log.CtxRequestID, rid)
		c.Header(log.HeaderRequestID, rid)
		c.Next()
	}
}
