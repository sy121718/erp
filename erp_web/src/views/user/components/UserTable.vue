<script setup lang="ts">
import type { UserInfo } from '@/api/userService'
import type { Pagination } from '../utils/types'

interface Props {
  data: UserInfo[]
  loading: boolean
  pagination: Pagination
}

defineProps<Props>()
const emit = defineEmits<{
  'edit': [row: UserInfo]
  'toggleStatus': [row: UserInfo]
  'resetPassword': [row: UserInfo]
  'pageChange': [page: number]
  'sizeChange': [size: number]
}>()

const columns = [
  { title: 'ID', key: 'id', width: 80, fixed: 'left' },
  { title: '用户名', key: 'username', width: 120 },
  { title: '姓名', key: 'name', width: 120 },
  { title: '邮箱', key: 'email', width: 180 },
  { title: '手机号', key: 'phone', width: 140 },
  { title: '积分', key: 'points', width: 100 },
  { title: '状态', key: 'status', width: 100, customSlot: 'status' },
  { title: '最后登录时间', key: 'last_login_time', width: 180 },
  { title: '注册时间', key: 'create_time', width: 180 },
  { title: '操作', key: 'action', width: 240, fixed: 'right', customSlot: 'action' }
]
</script>

<template>
  <lay-table
    :columns="columns"
    :dataSource="data"
    :loading="loading"
    :default-toolbar="true"
    :page="{
      current: pagination.page,
      limit: pagination.pageSize,
      total: pagination.total
    }"
    @change="(page: number) => emit('pageChange', page)"
    @limitChange="(size: number) => emit('sizeChange', size)"
  >
    <template #status="{ row }">
      <lay-tag
        :type="row.status === 1 ? 'normal' : 'danger'"
        style="cursor: pointer"
        @click="emit('toggleStatus', row)"
      >
        {{ row.status === 1 ? '正常' : '禁用' }}
      </lay-tag>
    </template>

    <template #action="{ row }">
      <lay-button size="xs" type="primary" @click="emit('edit', row)">编辑</lay-button>
      <lay-button size="xs" type="warm" @click="emit('resetPassword', row)">重置密码</lay-button>
      <lay-button
        size="xs"
        :type="row.status === 1 ? 'danger' : 'normal'"
        @click="emit('toggleStatus', row)"
      >
        {{ row.status === 1 ? '封禁' : '解封' }}
      </lay-button>
    </template>
  </lay-table>
</template>
