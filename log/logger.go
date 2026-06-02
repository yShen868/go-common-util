package log

import (
	"fmt"

	"github.com/yShen868/go-common-util/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

// L 返回根 Logger。初始化前为 no-op，避免空指针（业务应在启动时 Init）
func L() *zap.Logger {
	if l == nil {
		noop, _ := zap.NewProduction()
		return noop
	}
	return l
}

// Init 根据运行模式与日志级别创建全局 Logger（便捷方法，LogDir 默认 backup）
func Init(level, mode string) error {
	return InitWith(Options{Level: level, Mode: mode})
}

// InitFromAppConfig 根据应用配置初始化全局 Logger。
// 日志级别取自 cfg.Log.Level；输出方式由 cfg.Server.Mode 决定（debug→控制台，release→按日文件）。
func InitFromAppConfig(cfg *config.AppConfig) error {
	if cfg == nil {
		return fmt.Errorf("app config is nil")
	}
	return InitWith(Options{
		Level:  cfg.Log.Level,
		Mode:   cfg.Server.Mode,
		LogDir: cfg.Log.LogDir(),
	})
}

// InitWith 根据 Options 创建全局 Logger
func InitWith(opts Options) (err error) {
	var lvl zap.AtomicLevel
	if e := lvl.UnmarshalText([]byte(opts.Level)); e != nil {
		_ = lvl.UnmarshalText([]byte("info"))
	}

	if opts.Mode == "debug" {
		l, err = buildDebugLogger(lvl)
	} else {
		l, err = buildProductionLogger(lvl, opts.logDir())
	}
	if err != nil {
		return err
	}
	return nil
}

func buildDebugLogger(lvl zap.AtomicLevel) (*zap.Logger, error) {
	zcfg := zap.NewProductionConfig()
	zcfg.Level = lvl
	zcfg.Encoding = "dev_console"
	zcfg.EncoderConfig.EncodeTime = customTimeEncoder
	zcfg.OutputPaths = []string{"stdout"}
	zcfg.ErrorOutputPaths = []string{"stderr"}
	z, err := zcfg.Build(zap.AddCaller())
	if err != nil {
		return nil, fmt.Errorf("zap build debug: %w", err)
	}
	return z, nil
}

func buildProductionLogger(lvl zap.AtomicLevel, logDir string) (*zap.Logger, error) {
	encCfg := zap.NewProductionConfig().EncoderConfig
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = customTimeEncoder

	ws := zapcore.AddSync(newDailyFileWriter(logDir))
	core := zapcore.NewCore(
		newOrderedJSONEncoder(encCfg),
		ws,
		lvl,
	)
	return zap.New(core, zap.AddCaller()), nil
}
