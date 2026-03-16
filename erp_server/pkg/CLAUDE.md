# pkg 公共库

> 目录：`pkg/`

## config 配置

```go
config.Init("dev")
cfg := config.Get()
```

## database 数据库

```go
database.Init(&cfg.Database)
db := database.Get()
database.Close()
```

## errors 错误

```go
// 预定义错误
errors.ErrInvalidParams    // 400 参数错误
errors.ErrUnauthorized     // 401 未授权
errors.ErrForbidden        // 403 禁止访问
errors.ErrNotFound         // 404 不存在

// 创建错误
errors.NewBadRequest(code, "消息")
errors.NewUnauthorized(code, "消息")
errors.NewNotFound(code, "消息")
```

## log 日志

```go
log.Init(&cfg.Log)
log.Info("消息", zap.String("key", "value"))
log.Error("错误", zap.Error(err))
log.Sync()
```

## utils 响应

```go
utils.Success(c, data)
utils.Success(c, &utils.PaginationResponse{
    List: list, Total: total, Page: page, PageSize: pageSize,
})
```

## 目录结构

```
pkg/
├── config/      # 配置读取
├── database/    # 数据库连接
├── errors/      # 错误定义
├── log/         # 日志封装
├── utils/       # 工具函数
├── jwt/         # JWT
├── redis/       # Redis
├── captcha/     # 验证码
└── sign/        # 签名验证
```