package ginutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yShen868/go-common-util/config"
	"github.com/yShen868/go-common-util/log"
	"github.com/yShen868/go-common-util/response"
	"go.uber.org/zap"
)

// RegisterRoutes 在已创建并挂载默认中间件的 Engine 上注册业务路由
type RegisterRoutes func(e *gin.Engine)

// Run 加载配置、初始化日志、创建带默认中间件的 Engine、注册路由并阻塞监听 HTTP
func Run(register RegisterRoutes) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if err := log.InitFromAppConfig(cfg); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	gin.SetMode(cfg.Server.Mode)

	e := NewEngine()
	e.GET("/api/v1/health", response.Health)
	if register != nil {
		register(e)
	}

	addr := listenAddr(cfg)
	log.L().Info("http server starting", zap.String("addr", addr), zap.String("mode", cfg.Server.Mode))
	if err := e.Run(addr); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}

func listenAddr(cfg *config.AppConfig) string {
	if cfg == nil || cfg.Server.Port == 0 {
		return ":8080"
	}
	return fmt.Sprintf(":%d", cfg.Server.Port)
}
