# 公共组件规范

## 目录结构

```
src/components/
├── CLAUDE.md
├── SearchForm.vue      # 搜索表单
├── PageHeader.vue      # 页面头部
└── ...
```

## 组件规范

### Props 定义

```typescript
interface Props {
  title: string        // 必填
  visible?: boolean    // 可选
  data?: User[]        // 可选数组
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  data: () => []
})
```

### Emits 定义

```typescript
const emit = defineEmits<{
  update: [value: string]
  close: []
  change: [data: object]
}>()

emit('update', 'value')
```

### 样式规范

- 使用 CSS 变量 `var(--erp-*)`
- 遵循设计规范的间距、圆角、阴影

```vue
<style scoped>
.component {
  background-color: var(--erp-bg-card);
  border-radius: var(--erp-radius-lg);
  padding: 16px; /* p-4 */
}
</style>
```

## 注意事项

1. 公共组件应该是可复用的、通用的
2. 页面专用组件放在 `views/${页面}/components/`
3. 弹窗组件也放在页面 components 目录