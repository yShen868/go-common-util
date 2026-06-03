package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/log"
)

// Logger 将带 request_id 等字段的 Logger 注入 request context，业务层用 log.FromContext 取用
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.WithGin(c)
		ctx := log.NewContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
