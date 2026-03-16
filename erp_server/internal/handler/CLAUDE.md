# Handler 接口层

> 目录：`internal/handler/`

## 职责

- HTTP 请求处理
- 参数校验（binding 标签）
- DTO 定义和转换
- 响应返回

## Gin 风格要点

```go
// DTO 定义 - 使用 binding 标签
type LoginRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
}

// Handler 实现
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if !middleware.BindJSON(c, &req) {
        return  // 错误已由中间件处理
    }

    entity, err := h.svc.Login(c.Request.Context(), &req)
    if err != nil {
        c.Error(err)  // 记录错误，由中间件统一响应
        return
    }

    utils.Success(c, ToResponse(entity))
}
```

## 辅助函数

```go
middleware.BindJSON(c, &req)  // 绑定 JSON，失败自动记录错误
middleware.BindQuery(c, &req) // 绑定 Query
middleware.ParseID(c)         // 解析路径参数 ID
middleware.GetUserID(c)       // 获取当前用户 ID
```

## 中间件注册顺序

```go
r.Use(gin.Recovery())
r.Use(middleware.Logger())
r.Use(middleware.CORS())
r.Use(middleware.ErrorHandler())  // 最内层
```

## 目录结构

```
internal/handler/
├── middleware/
│   ├── error.go      # 错误处理
│   ├── auth.go       # 认证
│   ├── cors.go       # 跨域
│   └── logger.go     # 日志
└── admin/
    ├── handler.go    # HTTP 处理器
    ├── dto.go        # Request/Response DTO
    └── admin_router.go # 路由注册
```