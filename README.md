# go-common-util

个人 Go 工具库。

## 模块

| 包 | 说明 |
|---|---|
| `log` | 基于 zap 的日志初始化；支持开发控制台彩色输出与生产按日切片文件 |

## 使用

```go
import "github.com/yShen868/go-common-util/log"

if err := log.Init(log.Options{
    Level:  "info",
    Mode:   "release", // debug | release
    LogDir: "backup",
}); err != nil {
    // handle
}
log.L().Info("hello")
```

仓库：`git@github.com:yShen868/go-common-util.git`
