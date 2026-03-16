# ERP 数据采集推送平台

## 项目介绍

这是一个企业级 ERP 系统，用于商品数据采集和跨平台推送。

### 核心功能

- **数据采集**：从 1688 等平台采集商品数据
- **AI 优化**：使用 AI 优化商品文案和图片
- **多平台推送**：支持推送到亚马逊、拼多多、TikTok、Shopee、Lazada 等平台
- **多店铺管理**：支持多店铺绑定和管理
- **全局查找**：必须使用对应的lsp插件
---

## 项目结构

```
erp/
├── erp_web/                # 前端项目 (Vue 3)
│   ├── src/               # 源代码
│   └── CLAUDE.md          # 前端开发规范
│
├── erp_server/            # 后端项目 (Go)
│   ├── cmd/               # 应用入口
│   ├── configs/           # 配置文件
│   ├── internal/          # 私有代码 (DDD 分层)
│   │   ├── domain/        # 领域层
│   │   ├── application/    # 应用层
│   │   ├── infrastructure/# 基础设施层
│   │   └── handler/       # 接口层
│   ├── pkg/               # 公共库
│   ├── test/              # 测试代码
│   └── CLAUDE.md          # 后端开发规范
│
├── .cursor/rules/         # Cursor AI 规则
└── CLAUDE.md              # 项目总规则
```

---

## 技术栈

### 前端

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 框架 |
| Vite | - | 构建工具 |
| Layui Vue | - | UI 组件库 |
| Pinia | - | 状态管理 |
| Vue Router | - | 路由 |

### 后端

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.25+ | 语言 |
| Gin | - | Web 框架 |
| GORM | - | ORM |
| Viper | - | 配置管理 |
| Zap | - | 日志 |
| MySQL | - | 数据库 |

---


## 开发规范

### 前端开发

- 规则文件：`erp_web/CLAUDE.md`
- 使用 `<script setup>` 语法
- Layui Vue 组件直接使用

### 后端开发

- 规则文件：`erp_server/CLAUDE.md`
- DDD 分层架构
- 各层详细规则：
  - Domain 层：`erp_server/internal/domain/CLAUDE.md`
  - Application 层：`erp_server/internal/application/CLAUDE.md`
  - Infrastructure 层：`erp_server/internal/infrastructure/CLAUDE.md`
  - Handler 层：`erp_server/internal/handler/CLAUDE.md`
  - 公共库：`erp_server/pkg/CLAUDE.md`

---

## 快速开始

### 前端

```bash
cd erp_web
npm install
npm run dev
```

### 后端

```bash
cd erp_server
go mod tidy
go run cmd/api/main.go
```

服务启动后访问：`http://localhost:8081/health`

---

## 配置

### 后端配置

配置文件位于 `erp_server/configs/`：

- `config.yaml` - 主配置
- `config.dev.yaml` - 开发环境
- `config.prod.yaml` - 生产环境

### 数据库配置

```yaml
database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  username: root
  password: GKXTGztEWwkWtrWA
  dbname: erp
```

---
