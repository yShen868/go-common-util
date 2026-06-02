package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadOptions 控制配置文件加载行为。
type LoadOptions struct {
	// Env 运行环境，对应 config.{env}.yaml；空则读取环境变量 APP_ENV（默认 dev）
	Env string
	// ConfigPaths 配置文件搜索路径，默认 ["config", "."]
	ConfigPaths []string
}

// Load 从 config/ 目录加载 config.{dev|prod}.yaml。
// 环境由 APP_ENV 决定：prod 加载 config.prod.yaml，其余加载 config.dev.yaml。
func Load(opts ...LoadOptions) (*AppConfig, error) {
	var o LoadOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	paths := o.ConfigPaths
	if len(paths) == 0 {
		paths = []string{"config", "."}
	}

	v := viper.New()
	v.SetConfigName(configFileName(o.Env))
	v.SetConfigType("yaml")
	for _, p := range paths {
		v.AddConfigPath(p)
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &c, nil
}

func configFileName(env string) string {
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	switch env {
	case "prod":
		return "config.prod"
	default:
		return "config.dev"
	}
}
