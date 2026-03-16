<script setup lang="ts">
import { ref, computed } from 'vue'

interface TableColumn {
  title: string
  key?: string
  width?: string
  align?: 'left' | 'center' | 'right'
  fixed?: 'left' | 'right'
  slot?: string
  type?: 'checkbox' | 'radio' | 'index' | 'expand'
}

interface TablePagination {
  page: number
  limit: number
  total: number
}

interface Props {
  columns: TableColumn[]
  dataSource: any[]
  loading?: boolean
  pagination?: TablePagination
  showSelection?: boolean
  showIndex?: boolean
  rowKey?: string
  emptyText?: string
  height?: string | number
  stripe?: boolean
  border?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  showSelection: false,
  showIndex: false,
  rowKey: 'id',
  emptyText: '暂无数据',
  stripe: true,
  border: false
})

const emit = defineEmits<{
  selectionChange: [selectedKeys: (string | number)[]]
  pageChange: [page: number]
  limitChange: [limit: number]
  rowClick: [row: any]
}>()

// 选中行
const selectedKeys = ref<(string | number)[]>([])
const isAllSelected = computed(() => {
  if (props.dataSource.length === 0) return false
  return props.dataSource.every(row => selectedKeys.value.includes(row[props.rowKey]))
})

// 处理表格选择变化
const handleSelectionChange = (keys: any) => {
  selectedKeys.value = Array.isArray(keys) ? keys : [keys]
  emit('selectionChange', selectedKeys.value)
}

// 分页变化
const handlePageChange = (page: number) => {
  emit('pageChange', page)
}

// 每页条数变化
const handleLimitChange = (limit: number) => {
  emit('limitChange', limit)
}

// 行点击
const handleRowClick = (row: any) => {
  emit('rowClick', row)
}

// 计算序号
const getIndex = (index: number) => {
  if (!props.pagination) return index + 1
  return (props.pagination.page - 1) * props.pagination.limit + index + 1
}
</script>

<template>
  <div class="erp-table-wrapper">
    <!-- 表格 -->
    <lay-table
      :columns="columns as any"
      :dataSource="dataSource"
      :loading="loading"
      :height="height"
      :stripe="stripe"
      :border="border"
      :empty-text="emptyText"
      :row-key="rowKey"
      @selection-change="handleSelectionChange"
      @row-click="handleRowClick"
    >
      <!-- 序号列 -->
      <template v-if="showIndex" #index="{ index }">
        {{ getIndex(index) }}
      </template>

      <!-- 自定义列插槽 -->
      <template v-for="col in columns.filter(c => c.slot)" :key="col.key" #[col.slot!]="{ row, index }">
        <slot :name="col.slot" :row="row" :index="getIndex(index) - 1" />
      </template>
    </lay-table>

    <!-- 分页 -->
    <div v-if="pagination" class="erp-table-pagination">
      <div class="pagination-left">
        <slot name="pagination-left" />
      </div>
      <div class="pagination-right">
        <lay-page
          v-model="pagination.page"
          :total="pagination.total"
          :limit="pagination.limit"
          @change="handlePageChange"
          @limit-change="handleLimitChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.erp-table-wrapper {
  background-color: var(--erp-bg-card);
  border-radius: var(--erp-radius-lg);
  overflow: hidden;
}

/* 分页 */
.erp-table-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-top: 1px solid var(--erp-border-color);
}

.pagination-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 覆盖 Layui Vue 表格样式 */
:deep(.layui-table) {
  font-size: 13px;
}

:deep(.layui-table thead th) {
  background-color: #f9fafb;
  color: var(--erp-text-secondary);
  font-weight: 500;
  padding: 12px 16px;
}

:deep(.layui-table tbody td) {
  padding: 16px;
  color: var(--erp-text-primary);
}

:deep(.layui-table tbody tr:hover) {
  background-color: #f9fafb;
}

:deep(.layui-table--stripe .layui-table tbody tr:nth-child(even)) {
  background-color: rgba(243, 244, 246, 0.5);
}

/* 分页样式覆盖 */
:deep(.layui-page) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.layui-page .layui-this) {
  background-color: var(--erp-primary);
  border-color: var(--erp-primary);
}
</style>