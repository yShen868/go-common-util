package response

import "github.com/gin-gonic/gin"

// Health 简单健康检查，返回 {"status":"ok"}
func Health(c *gin.Context) {
	OK(c, gin.H{"status": "ok"})
}
