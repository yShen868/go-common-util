# go-common-util

个人 Go 工具库（多模块仓库）。

## 结构

```
go-common-util/
  go.work          # 本地工作区，聚合子模块
  config/
    go.mod         # module github.com/yShen868/go-common-util/config
    ...
  log/
    go.mod         # module github.com/yShen868/go-common-util/log
    ...
  response/
    go.mod         # module github.com/yShen868/go-common-util/response
    errno/         # 业务错误码
    ...
  ginutil/
    go.mod         # module github.com/yShen868/go-common-util/ginutil
    ...
```

| 模块路径 | 说明 |
|---|---|
| `github.com/yShen868/go-common-util/config` | YAML 应用配置结构体与 viper 加载（server/mysql/jwt/log） |
| `github.com/yShen868/go-common-util/log` | 基于 zap 的日志初始化；dev 控制台 / prod 按日文件；`FromContext` / `WithGin` |
| `github.com/yShen868/go-common-util/response` | 统一 JSON 响应、`errno` 子包、`Health` 处理器 |
| `github.com/yShen868/go-common-util/ginutil` | Gin 默认中间件链、`Run` 一键启动 HTTP 服务 |

根目录**没有** `go.mod`，只有 `go.work`；每个子目录（如 `log/`）是独立可发布的模块。

## 使用

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/yShen868/go-common-util/ginutil"
    "github.com/yShen868/go-common-util/response"
)

func main() {
    _ = ginutil.Run(func(e *gin.Engine) {
        e.GET("/api/v1/health", response.Health)
    })
}
```

## 发布版本

多模块仓库打 tag 需带子目录前缀，例如：

可用
```bash
log
git tag log/v0.4.0
git push origin log/v0.4.0

config
git tag config/v0.1.0
git push origin config/v0.1.0

```

```bash
git tag log/v0.1.0
git push origin log/v0.1.0
```

一次推所有 tag
```bash
git push origin --tags
```


引用：

```bash
go get github.com/yShen868/go-common-util/log@log/v0.1.0
```

仓库：`git@github.com:yShen868/go-common-util.git`


删除 tag
```bash
git tag -d v0.2.0    
git push origin :refs/tags/v0.2.0
```