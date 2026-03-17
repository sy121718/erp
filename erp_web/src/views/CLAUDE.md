# 页面视图规范

## 页面结构（参考 `admin/`）

```
page-name/
├── index.vue               # 页面编排：引用组件 + 解构 hook，禁止写业务逻辑
├── components/             # 页面专用组件
│   ├── XxxTable.vue        # 表格：columns + 插槽 + emit 事件
│   ├── XxxAddForm.vue      # 新增弹窗：自管 formData、校验、提交
│   └── XxxEditForm.vue     # 编辑弹窗：自加载详情、自管 formData、提交
└── utils/
    ├── types.ts            # 页面级类型（SearchForm、Pagination）
    └── hook.tsx            # 页面状态：搜索/分页/列表数据 + 弹窗开关 + 行操作
```

## 组件化规则

**必须拆为独立组件**：表格、新增弹窗、编辑弹窗、详情弹窗/抽屉

**禁止在 index.vue 中写**：`lay-layer`、`lay-table` + 列配置、表单校验/提交逻辑

## 各文件职责

| 文件 | 职责 | 不做 |
|------|------|------|
| `index.vue` | 引用组件、解构 hook、模板编排 | 业务逻辑、弹窗/表格模板 |
| `hook.tsx` | 搜索/分页/数据获取、弹窗 visible 控制、confirm 类行操作 | 表单数据、列配置 |
| `types.ts` | 页面级共享类型 | 业务实体类型（放 `api/`） |
| `XxxTable.vue` | props: data/loading/pagination，emit: edit/delete/pageChange 等 | 调 API |
| `XxxAddForm.vue` | props: v-model:visible，emit: success | 依赖外部状态 |
| `XxxEditForm.vue` | props: v-model:visible + targetId，emit: success | 依赖外部状态 |

## 弹窗统一结构

```vue
<lay-layer :modelValue="visible" @update:modelValue="(v: boolean) => emit('update:visible', v)"
  title="xxx" :area="['480px', 'auto']">
  <div class="dialog-form">
    <lay-form>...</lay-form>
    <div class="dialog-form__footer">
      <lay-button @click="emit('update:visible', false)">取消</lay-button>
      <lay-button type="primary" :loading="loading" @click="handleSubmit">确定</lay-button>
    </div>
  </div>
</lay-layer>
```

```css
.dialog-form { padding: 20px 24px 0; }
.dialog-form__footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding: 16px 0; border-top: 1px solid var(--erp-border-color); margin-top: 8px;
}
```

## 样式

- 页面容器 `padding: 16px`，卡片间距 `mb-4`，弹窗宽度 `480px`
- 使用 CSS 变量 `var(--erp-*)`
