<script setup lang="ts">
import { ref, onMounted } from 'vue'

const loading = ref(false)
const statistics = ref({
  totalProducts: 0,
  todayCollected: 0,
  pendingPush: 0,
  successPush: 0
})

const fetchStatistics = async () => {
  loading.value = true
  // TODO: 调用 API 获取统计数据
  loading.value = false
}

onMounted(() => {
  fetchStatistics()
})
</script>

<template>
  <div class="page-container">
    <div class="statistics-row">
      <lay-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ statistics.totalProducts }}</div>
          <div class="stat-label">商品总数</div>
        </div>
      </lay-card>
      <lay-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ statistics.todayCollected }}</div>
          <div class="stat-label">今日采集</div>
        </div>
      </lay-card>
      <lay-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ statistics.pendingPush }}</div>
          <div class="stat-label">待推送</div>
        </div>
      </lay-card>
      <lay-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ statistics.successPush }}</div>
          <div class="stat-label">推送成功</div>
        </div>
      </lay-card>
    </div>

    <lay-card class="welcome-card">
      <div class="welcome-content">
        <h2>欢迎使用妙手ERP</h2>
        <p>跨境电商数据采集与多平台推送管理系统</p>
      </div>
    </lay-card>
  </div>
</template>

<style scoped>
.page-container {
  padding: 16px;
}

.statistics-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.stat-content {
  padding: 16px 0;
}

.stat-value {
  font-size: 32px;
  font-weight: 600;
  color: var(--erp-primary);
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: var(--erp-text-secondary);
}

.welcome-card {
  margin-top: 16px;
}

.welcome-content {
  text-align: center;
  padding: 48px 0;
}

.welcome-content h2 {
  font-size: 24px;
  color: var(--erp-text-primary);
  margin-bottom: 12px;
}

.welcome-content p {
  font-size: 14px;
  color: var(--erp-text-secondary);
}
</style>
