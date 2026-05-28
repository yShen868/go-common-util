package log

// Options 日志初始化参数。
type Options struct {
	// Level: debug, info, warn, error
	Level string
	// Mode: debug 输出到 stdout（dev_console）；release 写入 LogDir/YYYYMMDD.log（ordered_json）
	Mode string
	// LogDir 生产模式日志目录，默认 backup
	LogDir string
}

func (o Options) logDir() string {
	if o.LogDir == "" {
		return "backup"
	}
	return o.LogDir
}
