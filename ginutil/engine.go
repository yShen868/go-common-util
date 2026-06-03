package ginutil

import "github.com/gin-gonic/gin"

// NewEngine 创建 Gin 引擎并挂载默认中间件链（调用方在 Run 前需已 gin.SetMode）
func NewEngine() *gin.Engine {
	e := gin.New()
	e.Use(
		RequestID(),
		Logger(),
		RecoveryZap(),
		TimingStats(),
		AccessLog(),
	)
	return e
}
