<script setup lang="ts">
import { ref, onMounted } from 'vue'

const loading = ref(false)
const tableData = ref([])

const fetchData = async () => {
  loading.value = true
  // TODO: 调用 API 获取数据
  loading.value = false
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="page-container">
    <lay-card>
      <template #title>
        <span class="card-title">公用采集箱</span>
      </template>
      <lay-table :columns="[]" :dataSource="tableData" :loading="loading">
        <template #empty>
          <lay-empty description="暂无数据" />
        </template>
      </lay-table>
    </lay-card>
  </div>
</template>

<style scoped>
.page-container {
  padding: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--erp-text-primary);
}
</style>
