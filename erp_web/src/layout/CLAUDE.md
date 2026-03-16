# 布局组件规范

## 目录结构

```
src/layout/
├── CLAUDE.md
├── index.vue                    # 主布局入口
└── components/
    ├── lay-sidebar/
    │   ├── NavVertical.vue      # 垂直侧边栏
    │   └── components/
    │       ├── SidebarItem.vue
    │       ├── SidebarLogo.vue
    │       ├── SidebarLinkItem.vue
    │       └── SidebarBreadCrumb.vue
    ├── lay-navbar/
    │   └── index.vue            # 顶部导航栏
    └── lay-content/
        └── index.vue            # 内容区域
```

## 布局规格

| 区域 | 规格 |
|------|------|
| 顶部导航 | 48px 高，背景 `var(--erp-primary)` |
| 左侧边栏 | 192px 宽，白色背景，右边框 |
| 主内容区 | 灰色背景 `var(--erp-bg-base)`，16px 内边距 |

## 组件职责

### NavVertical (侧边栏)

- 显示 Logo
- 渲染菜单列表
- 支持折叠状态
- Props: `menuList`, `collapsed`, `logoTitle`, `activePath`

### LayNavbar (顶部导航)

- 折叠/展开按钮
- 面包屑导航
- 用户信息下拉菜单
- Props: `collapsed`, `userInfo`, `breadcrumbs`

### LayContent (内容区)

- 包裹 router-view
- 提供统一内边距
- Props: `padding`, `background`

## 样式规范

```vue
<style scoped>
/* 使用 CSS 变量 */
.nav-vertical {
  width: 192px;
  background-color: var(--erp-bg-card);
  border-right: 1px solid var(--erp-border-color);
}
</style>
```