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
```

| 模块路径 | 说明 |
|---|---|
| `github.com/yShen868/go-common-util/config` | YAML 应用配置结构体与 viper 加载（server/mysql/jwt/log） |
| `github.com/yShen868/go-common-util/log` | 基于 zap 的日志初始化；dev 控制台 / prod 按日文件；支持 `InitFromAppConfig` |

根目录**没有** `go.mod`，只有 `go.work`；每个子目录（如 `log/`）是独立可发布的模块。

## 使用

```go
import (
    "github.com/yShen868/go-common-util/config"
    "github.com/yShen868/go-common-util/log"
)

cfg, err := config.Load()
if err != nil {
    // handle
}
if err := log.InitFromAppConfig(cfg); err != nil {
    // handle
}
log.L().Info("hello")
```

## 发布版本

多模块仓库打 tag 需带子目录前缀，例如：

可用
```bash
git tag log/v0.3.0
git push origin log/v0.3.0
```

```bash
git tag log/v0.1.0
git push origin log/v0.1.0
```


```bash
git tag log/v0.3.0 v0.3.0
git push origin log/v0.3.0
```


引用：

```bash
go get github.com/yShen868/go-common-util/log@log/v0.1.0
```

仓库：`git@github.com:yShen868/go-common-util.git`

