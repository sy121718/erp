# Go + DDD 开发规范

## 核心规则

- 用中文回答，执行前先提示
- 代码完成后必须检查逻辑漏洞和错误
- 文件删除前需要询问确认

## 项目信息

| 组件 | 版本/说明 |
|------|----------|
| Go | 1.25+ |
| 架构 | DDD（领域驱动设计） |
| Web 框架 | Gin |
| ORM | GORM |
| 配置管理 | Viper |
| 日志 | Zap |

## Gin 官方风格

> 参考：https://gin-gonic.com/zh-cn/docs/

### 核心要点

```go
// 路由器
r := gin.Default()  // 默认中间件
r := gin.New()      // 无中间件

// 绑定（推荐 ShouldBind 系列）
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}

// 响应
c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data})
```

### 绑定方法

| 方法 | 说明 |
|------|------|
| `ShouldBind`, `ShouldBindJSON`, `ShouldBindQuery` | 推荐，失败返回 error |
| `Bind`, `BindJSON`, `BindQuery` | 不推荐，失败自动终止请求 |

---

## DDD 分层架构

```
Handler 层：Request DTO → Entity → Response DTO
    ↓
Application 层：Entity（业务逻辑）
    ↓
Infrastructure 层：Entity → PO → 数据库
    ↓
Domain 层：Entity（充血模型、值对象、仓储接口）
```

### 各层规则文件

| 层 | 文件 | 职责 |
|---|------|------|
| Domain | `internal/domain/CLAUDE.md` | 实体、值对象、仓储接口 |
| Application | `internal/application/CLAUDE.md` | 业务流程编排 |
| Infrastructure | `internal/infrastructure/CLAUDE.md` | Repository 实现、PO 转换 |
| Handler | `internal/handler/CLAUDE.md` | HTTP 处理、DTO |
| pkg | `pkg/CLAUDE.md` | 公共库 |

---

## 错误处理

**使用 `c.Error()` + 中间件统一处理**

```go
// Handler 层
if err != nil {
    c.Error(err)
    return
}

// 辅助函数
middleware.BindJSON(c, &req)  // 绑定 JSON
middleware.BindQuery(c, &req) // 绑定 Query
middleware.ParseID(c)         // 解析 ID
```

---

## 响应格式

```json
// 成功
{"code": 0, "msg": "success", "data": {...}}

// 分页
{"code": 0, "msg": "success", "data": {"list": [...], "total": 100, "page": 1, "page_size": 10}}

// 错误
{"code": 400, "msg": "参数错误"}
```

### 错误码

| 码 | 说明 |
|---|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |

---

## 路由规范

**只使用 GET 和 POST**

| 操作 | 方法 | URL |
|------|------|-----|
| 列表 | GET | `/api/admin/list` |
| 详情 | GET | `/api/admin/:id` |
| 创建 | POST | `/api/admin/create` |
| 更新 | POST | `/api/admin/update/:id` |
| 删除 | POST | `/api/admin/delete/:id` |

---

## 开发流程

```
编写接口 → 单元测试 → 集成测试 → 同步 Apipost → 下一个接口
```

**禁止一次性生成多个接口**

---

## 数据库字段

| 字段 | 类型 |
|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT |
| create_time, update_time, deleted_time | DATETIME(3) |
| create_by, update_by | BIGINT UNSIGNED |

---

## 目录结构

```
erp_server/
├── cmd/api/main.go
├── configs/
├── internal/
│   ├── domain/          # 领域层
│   ├── application/     # 应用层
│   ├── infrastructure/  # 基础设施层
│   └── handler/         # 接口层
├── pkg/                 # 公共库
└── test/                # 测试
```