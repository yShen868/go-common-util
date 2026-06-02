package config

// AppConfig 应用配置根结构，与 config/*.yaml 顶层字段一一对应。
type AppConfig struct {
	// Server HTTP 服务相关配置
	Server ServerConfig `mapstructure:"server"`
	// MySQL 数据库连接池配置
	MySQL MySQLConfig `mapstructure:"mysql"`
	// JWT 鉴权令牌配置
	JWT JWTConfig `mapstructure:"jwt"`
	// Log 日志输出配置（级别、目录等）
	Log LogConfig `mapstructure:"log"`
}

// ServerConfig HTTP 服务配置。
type ServerConfig struct {
	// Port 监听端口；为 0 时业务层可回退为 8080
	Port int `mapstructure:"port"`
	// Mode Gin 运行模式：debug（开发，彩色控制台日志）| release（生产，按日写文件）
	Mode string `mapstructure:"mode"`
}

// MySQLConfig MySQL 数据源与连接池配置。
type MySQLConfig struct {
	// DSN 数据源名称，格式：user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	DSN string `mapstructure:"dsn"`
	// MaxIdleConns 连接池中最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// MaxOpenConns 连接池中最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
}

// JWTConfig JWT 签发与校验配置。
type JWTConfig struct {
	// Secret 签名密钥，生产环境务必替换为强随机值
	Secret string `mapstructure:"secret"`
	// Expire 令牌有效期（秒），例如 7200 表示 2 小时
	Expire int `mapstructure:"expire"`
}

// LogConfig 日志模块配置。
type LogConfig struct {
	// Level 日志级别：debug | info | warn | error；dev 建议 debug，prod 建议 info
	Level string `mapstructure:"level"`
	// Dir 生产模式（server.mode=release）日志目录，空则使用默认 backup
	Dir string `mapstructure:"dir"`
}

// LogDir 返回生产模式日志目录，未配置时为 backup。
func (c LogConfig) LogDir() string {
	if c.Dir == "" {
		return "backup"
	}
	return c.Dir
}
