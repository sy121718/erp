# 页面视图规范

## 目录结构

```
src/views/
├── CLAUDE.md
├── dashboard/
│   └── index.vue           # 首页
├── user/
│   ├── index.vue           # 用户列表
│   └── components/         # 用户页面专用组件
│       ├── UserForm.vue    # 用户表单弹窗
│       └── UserDetail.vue  # 用户详情抽屉
├── login-admin/            # 复杂页面示例
│   ├── index.vue           # 页面视图 (纯模板)
│   └── utils/              # 业务逻辑抽取
│       ├── types.ts        # 类型定义
│       └── hook.tsx        # 组合式函数 (useLogin)
└── product/
    ├── index.vue
    └── components/
        └── ...
```

## 页面结构

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userService } from '@/api'
import UserForm from './components/UserForm.vue'

// 列表数据
const loading = ref(false)
const tableData = ref([])
const pagination = ref({ page: 1, size: 10, total: 0 })

// 弹窗控制
const formVisible = ref(false)
const formRef = ref()

// 方法
const fetchData = async () => {
  loading.value = true
  const res = await userService.getList(pagination.value)
  tableData.value = res.list
  pagination.value.total = res.total
  loading.value = false
}

const handleAdd = () => {
  formVisible.value = true
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="page-container">
    <!-- 搜索区域 -->
    <lay-card class="mb-4">
      <!-- 搜索表单 -->
    </lay-card>

    <!-- 表格区域 -->
    <lay-card>
      <div class="flex justify-between mb-4">
        <lay-button type="primary" @click="handleAdd">新增</lay-button>
      </div>
      <lay-table :columns="columns" :dataSource="tableData" :loading="loading" />
      <lay-page v-model="pagination.page" :total="pagination.total" />
    </lay-card>

    <!-- 弹窗 -->
    <UserForm v-model:visible="formVisible" @success="fetchData" />
  </div>
</template>

<style scoped>
.page-container {
  padding: 16px;
}
</style>
```

## 页面规范

### 列表页面

- 搜索区域 + 表格区域 + 分页
- 操作按钮放在表格上方
- 使用 `lay-card` 包裹

### 表单弹窗

- 放在 `components/` 目录
- 使用 `v-model:visible` 控制显示
- 提交成功后 emit `success` 事件

### 样式

- 页面容器使用 16px 内边距
- 卡片间距使用 `mb-4` (16px)
- 使用 CSS 变量 `var(--erp-*)`

### 复杂页面拆分规范

当页面逻辑较复杂时，必须拆分为以下结构：

```
page-name/
├── index.vue       # 页面视图，只保留模板
└── utils/
    ├── types.ts    # 类型定义
    └── hook.tsx    # 组合式函数
```

**types.ts 规范：**
- 定义页面相关的所有 TypeScript 接口和类型
- 导出类型供 hook 和 index.vue 使用

**hook.tsx 规范：**
- 使用 `use` 前缀命名组合式函数，如 `useLogin`、`useUserList`
- 包含所有状态 (ref、reactive、computed)
- 包含所有方法 (事件处理、API 调用)
- 返回页面需要的所有数据和方法

**index.vue 规范：**
- 从 `./utils/hook` 导入组合式函数
- 只保留模板代码，禁止在 `<script>` 中写业务逻辑

**示例：**

```typescript
// utils/types.ts
export interface LoginForm {
  username: string
  password: string
}

// utils/hook.tsx
export const useLogin = () => {
  const formData = ref<LoginForm>({ username: '', password: '' })
  const loading = ref(false)
  
  const handleLogin = async () => { /* ... */ }
  
  return { formData, loading, handleLogin }
}

// index.vue
<script setup lang="ts">
import { useLogin } from './utils/hook'
const { formData, loading, handleLogin } = useLogin()
</script>
```