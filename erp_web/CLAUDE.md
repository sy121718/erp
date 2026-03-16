# ERP Web 前端开发规范

## 技术栈

Vue 3 + Vite + Layui Vue + Pinia + Vue Router + TypeScript

## 目录结构

```
src/
├── api/          # API 接口 → 详见 src/api/CLAUDE.md
├── components/   # 公共组件 → 详见 src/components/CLAUDE.md
├── layout/       # 框体组件 → 详见 src/layout/CLAUDE.md
├── composables/  # 组合式函数
├── router/       # 路由配置
├── stores/       # Pinia 状态管理 → 详见 src/stores/CLAUDE.md
├── utils/        # 工具函数
└── views/        # 页面视图 → 详见 src/views/CLAUDE.md
```

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 组件 | PascalCase | `UserList.vue` |
| 组合式函数 | use前缀 | `useUser.ts` |
| API | camelCase | `userService.ts` |
| 常量 | UPPER_SNAKE | `API_URL` |

---

## 设计规范

### 颜色系统

使用 CSS 变量，定义在 `src/style.css`：

| 分类 | 变量 | 用途 |
|------|------|------|
| 主色调 | `--erp-primary` (#03b3b2) | 按钮、链接、激活态 |
| 强调色 | `--erp-orange/green/red/blue/purple` | 特殊按钮、状态标签 |
| 中性色 | `--erp-bg-base/border-color` | 背景、边框 |
| 文字 | `--erp-text-primary/secondary/tertiary` | 主/次/辅助文字 |

### 布局规范

| 区域 | 规格 |
|------|------|
| 顶部导航 | 48px 高，背景 `#03b3b2`，文字白色 |
| 左侧边栏 | 192px 宽，白色背景，右边框 |
| 主内容区 | 灰色背景，16px 内边距 |
| 卡片/面板 | 白色背景，8px 圆角，shadow-sm |

### 间距规范

基于 4px 单位: xs(4px) / sm(8px) / md(12px) / lg(16px) / xl(24px)

常用: 卡片内边距 `p-4`, 表格单元格 `px-2 py-3`, 按钮组间距 `gap-2`

---

## 组件规范

### Layui Vue 主题覆盖

已在 `src/style.css` 中覆盖 Layui Vue 组件主题：

| 组件 | 覆盖内容 |
|------|----------|
| 按钮 | 主色调、危险色、成功色 |
| 菜单 | 激活态、悬停态 |
| 输入框 | 聚焦边框、阴影 |
| 标签 | success/danger/warning 颜色 |
| 表格 | 行悬停背景 |

### 组件开发规范

```vue
<script setup lang="ts">
// 1. 导入
// 2. Props / Emits (使用 TypeScript interface)
// 3. 状态与方法
</script>

<template>
  <!-- 使用 CSS 变量，避免硬编码颜色 -->
</template>

<style scoped>
/* 使用 var(--erp-*) 变量 */
</style>
```

---

## 子目录规范文件

| 目录 | 规范文件 | 内容 |
|------|----------|------|
| `src/api/` | CLAUDE.md | API 接口规范、请求封装 |
| `src/components/` | CLAUDE.md | 公共组件规范 |
| `src/layout/` | CLAUDE.md | 布局组件规范 |
| `src/stores/` | CLAUDE.md | Pinia Store 规范 |
| `src/views/` | CLAUDE.md | 页面视图规范 |

---

## 注意事项

1. 使用 `<script setup lang="ts">` 语法
2. Layui Vue 组件无需导入，直接使用
3. 使用 CSS 变量 `var(--erp-*)`，避免硬编码颜色
4. 复杂逻辑抽取为组合式函数
5. 开发前先阅读对应目录的 CLAUDE.md
6. **TypeScript 严格类型**：所有回调函数参数必须显式声明类型，禁止隐式 `any`
   - Axios 拦截器：`InternalAxiosRequestConfig`、`AxiosResponse`、`AxiosError<T>`
   - 示例：`(config: InternalAxiosRequestConfig) => { ... }`